"""
Main parser module for structured-parse.

This module provides the StructuredParser class that wraps the WebAssembly module
to parse labeled/structured text.
"""

import json
import os
import subprocess
import shutil
from dataclasses import dataclass, field
from typing import Any, Dict, List, Optional


@dataclass
class Label:
    """
    Label defines a label for parsing with options for required, dependencies, JSON, and block start.
    
    Attributes:
        name: Name of the label (case-insensitive matching, original casing preserved in results)
        required: Whether this label is required
        required_with: List of other label names required with this one
        is_json: Whether this label should be parsed as JSON
        is_block_start: Whether this label starts a new block
    """
    name: str
    required: bool = False
    required_with: List[str] = field(default_factory=list)
    is_json: bool = False
    is_block_start: bool = False

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for JSON serialization."""
        result = {"name": self.name}
        if self.required:
            result["required"] = self.required
        if self.required_with:
            result["requiredWith"] = self.required_with
        if self.is_json:
            result["isJson"] = self.is_json
        if self.is_block_start:
            result["isBlockStart"] = self.is_block_start
        return result


@dataclass
class ParserOptions:
    """
    ParserOptions allows customization of parser behavior.
    
    Attributes:
        separators: Allowed separator characters. Default is ":~-="
    """
    separators: Optional[str] = None

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for JSON serialization."""
        if self.separators:
            return {"separators": self.separators}
        return {}


@dataclass
class ParseResult:
    """
    ParseResult contains the parsed data and any errors encountered.
    
    Attributes:
        result: The parsed result map (label name -> value)
        errors: List of parsing/validation errors
    """
    result: Dict[str, Any]
    errors: List[str]


@dataclass
class ParseBlocksResult:
    """
    ParseBlocksResult contains the parsed blocks and any errors encountered.
    
    Attributes:
        blocks: Array of parsed blocks, each block is a map of label name -> value
        errors: List of parsing/validation errors
    """
    blocks: List[Dict[str, Any]]
    errors: List[str]


class StructuredParser:
    """
    StructuredParser provides WebAssembly-powered parsing of labeled/structured LLM output.
    
    This class wraps a WASM module compiled from Go to provide fast, consistent parsing
    across multiple programming languages.
    """

    def __init__(self, options: Optional[ParserOptions] = None, *, wasm_path: Optional[str] = None):
        """
        Creates a new StructuredParser instance.
        
        Args:
            options: Optional parser options (e.g., custom separators). If None, uses default options.
            wasm_path: Optional path to the WASM file (keyword-only). If not provided, uses the bundled WASM module.
            
        Raises:
            FileNotFoundError: If WASM module or wasmtime CLI is not found.
            
        Example:
            >>> parser = StructuredParser()  # Default options
            >>> parser = StructuredParser(ParserOptions(separators=":"))  # Custom separators
        """
        # Store options for later use
        self._default_options = options
        
        # Determine the path to the WASM file
        if wasm_path is None:
            # Use the bundled WASM module
            module_dir = os.path.dirname(__file__)
            wasm_path = os.path.join(module_dir, "wasm", "structured-parse.wasm")
        
        if not os.path.exists(wasm_path):
            raise FileNotFoundError(f"WASM module not found at: {wasm_path}")
        
        # Check if wasmtime CLI is available
        if not shutil.which("wasmtime"):
            raise FileNotFoundError(
                "wasmtime CLI not found. Please install wasmtime:\n"
                "  Linux/macOS: curl https://wasmtime.dev/install.sh -sSf | bash\n"
                "  macOS (Homebrew): brew install wasmtime\n"
                "  Ubuntu/Debian: sudo apt install wasmtime\n"
                "  Visit: https://wasmtime.dev/"
            )
        
        self.wasm_path = wasm_path

    def _call_wasm(self, command: str, labels: List[Label], text: str, options: Optional[ParserOptions] = None) -> Dict[str, Any]:
        """
        Internal method to call the WASM module.
        
        Args:
            command: The command to execute ("parse" or "parseBlocks")
            labels: Array of label definitions
            text: The text to parse
            options: Optional parser options
            
        Returns:
            The response from the WASM module
            
        Raises:
            RuntimeError: If the WASM call fails
        """
        # Prepare the request
        request = {
            "command": command,
            "labels": [label.to_dict() for label in labels],
            "text": text,
        }
        if options:
            request["options"] = options.to_dict()
        
        request_json = json.dumps(request)
        
        # Use wasmtime CLI to execute the WASM module
        try:
            result = subprocess.run(
                ["wasmtime", "run", self.wasm_path],
                input=request_json.encode("utf-8"),
                capture_output=True,
                check=True,
            )
            response_json = result.stdout.decode("utf-8").strip()
        except subprocess.CalledProcessError as e:
            raise RuntimeError(f"WASM execution failed: {e.stderr.decode('utf-8')}")
        
        # Parse the response
        try:
            response = json.loads(response_json)
        except json.JSONDecodeError as e:
            raise RuntimeError(f"Failed to parse WASM response: {e}\nResponse: {response_json}")
        
        # Check for system errors (not parsing errors)
        # If there's an "error" field, it's a system error (e.g., invalid request)
        # If ok is false but there are "errors" (plural), those are parsing errors and we return them
        if "error" in response:
            raise RuntimeError(f"WASM call failed: {response['error']}")
        
        return response

    def parse(self, labels: List[Label], text: str, options: Optional[ParserOptions] = None) -> ParseResult:
        """
        Parses text into a map of label names to their values.
        
        Args:
            labels: Array of label definitions
            text: The text to parse
            options: Optional parser options. If None, uses options from constructor.
            
        Returns:
            ParseResult containing the parsed data and any errors
            
        Example:
            >>> parser = StructuredParser()
            >>> labels = [
            ...     Label(name="Action", required=True),
            ...     Label(name="Parameters", is_json=True),
            ... ]
            >>> result = parser.parse(labels, "Action: test\\nParameters: {}")
            >>> print(result.result)
            {'Action': 'test', 'Parameters': {}}
        """
        # Use provided options or fall back to constructor options
        effective_options = options if options is not None else self._default_options
        response = self._call_wasm("parse", labels, text, effective_options)
        return ParseResult(
            result=response.get("result", {}),
            errors=response.get("errors", []),
        )

    def parse_blocks(self, labels: List[Label], text: str, options: Optional[ParserOptions] = None) -> ParseBlocksResult:
        """
        Parses text into blocks, splitting at the block start label.
        
        Args:
            labels: Array of label definitions (one must have is_block_start=True)
            text: The text to parse
            options: Optional parser options. If None, uses options from constructor.
            
        Returns:
            ParseBlocksResult containing the parsed blocks and any errors
            
        Example:
            >>> parser = StructuredParser()
            >>> labels = [
            ...     Label(name="Task", is_block_start=True),
            ...     Label(name="Result"),
            ... ]
            >>> result = parser.parse_blocks(labels, "Task: A\\nResult: 1\\nTask: B\\nResult: 2")
            >>> print(len(result.blocks))
            2
        """
        # Use provided options or fall back to constructor options
        effective_options = options if options is not None else self._default_options
        response = self._call_wasm("parseBlocks", labels, text, effective_options)
        return ParseBlocksResult(
            blocks=response.get("result", []),
            errors=response.get("errors", []),
        )

    def version(self) -> str:
        """
        Returns the version of the WASM module.
        
        Returns:
            Version string
        """
        response = self._call_wasm("version", [], "", None)
        return response.get("result", "unknown")



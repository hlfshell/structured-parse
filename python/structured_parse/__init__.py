"""
structured-parse: Multi-language library for parsing labeled/structured LLM output.

This library provides WebAssembly-powered parsing of labeled/structured text,
commonly used with LLM outputs. It supports block parsing, JSON fields,
required field validation, and robust handling of LLM output quirks.
"""

from .parser import StructuredParser, Label, ParserOptions, ParseResult, ParseBlocksResult
from .version import __version__

__all__ = [
    "StructuredParser",
    "Label",
    "ParserOptions",
    "ParseResult",
    "ParseBlocksResult",
    "__version__",
]



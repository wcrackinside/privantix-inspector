# Skills / Capabilities

## repository_scanning
Recursively traverses a target path and discovers supported files.

## file_format_detection
Maps file extensions to supported analyzers.

## encoding_detection
Detects BOM, validates UTF-8 and assigns a conservative fallback.

## delimiter_detection
Infers delimiter from sample rows using candidate scoring.

## header_detection
Uses lightweight heuristics to determine whether the first row is a header.

## column_profiling
Infers type, null rate, max length and sample values per column.

## metadata_extraction
Captures file metadata such as size, modified date and repository depth.

## rules_evaluation
Executes MVP hardcoded rules over file and column profiles.

## report_generation
Exports JSON, CSV and HTML outputs.

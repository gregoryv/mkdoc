# Changelog

## [0.4.0]

- <list of requirements> includes a parsed list of RFC 2119
  requirement sentences in the document
- Warn if empty lines only contain space characters
- Warn if adjecent sentences are Not separated by two space characters

## [0.3.0] 2025-05-20

- Warn if requirements are not tagged with (#R...)
- Rename <incfile ...> to <cat ...>; shorter and more memorable

## [0.2.1] 2025-05-19

- Fix links in beginning of line, where replaced with #ref-...

## [0.2.0] 2025-05-19

- Remove R prefix from requirement links
- Add parsing of requirement links
- Replace section links, e.g. (§1.1) with link to #section-1.1
- Fix references; lines starting with [\d+]\s are named and
  missing first character
  

## [0.1.0] 2025-05-17

Initial revision supporting 

- RFC like indentation
- table of contents generation
- reference link processing
- Go like reference links

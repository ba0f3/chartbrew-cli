---
name: cli-tool
description: Agent-first CLI tool template and examples.
---

# CLI Tool Skill

This skill defines the operational boundaries and instructions for using the CLI tool.

## Key Concepts

- Always use the `--output=json` or `--output=markdown` flags when extracting data, as this ensures you receive a highly structured payload.
- Do NOT try to bypass interactive prompts by using bash tools like `yes`. Use the documented stdin/env-var mechanisms.
- Read credentials natively using `--password -` and passing via stdin instead of putting passwords into command arguments directly.

## Standard Usage

1. **Information Gathering**: Run `cli help` and `cli <command> --help` to understand the available endpoints.
2. **Setup**: Use environment variables or non-interactive flags to configure the CLI.
3. **Execution**: Avoid running commands that could delete resources without explicitly checking constraints first.

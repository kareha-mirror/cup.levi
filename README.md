# levi - Least Enhanced vi

[Japanese version of README](README.ja.md)

**levi** is a text editor that aims to recreate the feel of classic **vi** on modern systems.

Its behavior is primarily inspired by **nvi**, although complete compatibility is not a goal.
Likewise, levi is **not** intended to become a feature-rich editor like **Vim**.

The design philosophy of levi is **Least Enhanced**: preserve the simplicity and responsiveness of traditional vi while adding only the features that have become essential in modern environments.

Accordingly, levi intentionally does **not** implement:

* Split windows
* Syntax highlighting
* LSP integration

Instead, levi provides only a few carefully chosen enhancements:

* Full support for UTF-8, including wide characters
* Access to the system clipboard
* Shared registers for copying text between levi processes
* A minimal color scheme system

> **Least Enhanced is not Minimal.**

levi is written in Go and can be built on any platform supported by Go.

## System Clipboard

Like many extended vi implementations, levi maps the `+` register to the system clipboard.

For example, `"+yy` copies the current line to the system clipboard, and `"+p` pastes the clipboard contents into the editor.

## Shared Registers

By default, the `x`, `y`, and `z` registers are shared between levi processes.

This makes it easy to copy text from one levi instance and paste it into another.

For example, use `"xyy` in one levi process to copy the current line into a shared register, then use `"xp` in another levi process to paste it.

Since levi intentionally does not provide split windows or a client-server architecture, shared registers offer a simple way to exchange text between multiple terminals.

## Minimal Color Schemes

levi includes a small collection of built-in color schemes.

Use

```text
:colors
```

to list the available color schemes.

Select one with

```text
:colors <name>
```

For example,

```text
:colors violet
```

selects the `violet` color scheme.

To make the selection persistent, specify it in the configuration file.

On most platforms, the configuration file is located at

```text
$HOME/.config/levi/editor.yaml
```

## Build

This project uses [Task](https://taskfile.dev/) as task runner.
To build levi, first install Task by:

```sh
go install github.com/go-task/task/v3/cmd/task@latest
```

Next, run `task` in project tree.

```sh
task
```

# Wii SYSCONF Tool

A simple CLI tool to convert Nintendo Wii `SYSCONF` files to and from JSON format.

This tool is intended for safe editing of Wii system configuration without hex editors.

---

## What it does

* Decode `SYSCONF` → JSON
* Encode JSON → valid `SYSCONF`

---

## Usage

### Decode

```bash
sysconfer decode SYSCONF
```

Creates:

```
SYSCONF.json
```

### Encode

```bash
sysconfer encode SYSCONF.json
```

Creates:

```
SYSCONF.json.bin
```

---

## JSON Format Rules

### Top level

The JSON file is a **dictionary of items**, where:

* **Key** = item name (UTF-8)
* **Value** = item definition

```json
{
  "BT.SENS": {
    "Type": "LONG",
    "hex": "00000003"
  }
}
```

---

## Item Name Rules

* Must be **UTF-8**
* Maximum length: **31 bytes**
* Must be unique

---

## Supported Types

| Type         | Size     |
| ------------ | -------- |
| `BIGARRAY`   | variable |
| `SMALLARRAY` | variable |
| `BYTE`       | 1 byte   |
| `SHORT`      | 2 bytes  |
| `LONG`       | 4 bytes  |
| `LONGLONG`   | 8 bytes  |
| `BOOL`       | 1 byte   |

---

##  Notes

* Lookup table is filled with zeroes. The System Menu will fall back to name-based lookup automatically
* Invalid or unsupported data will cause encoding to fail

---

## References

* WiiBrew SYSCONF documentation
  [https://wiibrew.org/wiki//shared2/sys/SYSCONF](https://wiibrew.org/wiki//shared2/sys/SYSCONF)

---

MIT Licensed

#!/usr/bin/env python3
"""Fail CI if Android backup rules start exporting application databases."""

from pathlib import Path
import xml.etree.ElementTree as ET


ROOT = Path(__file__).resolve().parents[1]
ANDROID = ROOT / "clients" / "android" / "app" / "src" / "main"
SENTINEL = "a2m_backup_disabled_sentinel"


def allowed_files(container: ET.Element) -> list[tuple[str, str]]:
    return [
        (item.attrib.get("domain", ""), item.attrib.get("path", ""))
        for item in container.findall("include")
    ]


legacy = ET.parse(ANDROID / "res/xml/backup_rules.xml").getroot()
assert legacy.tag == "full-backup-content", "Unexpected legacy backup rules root"
assert allowed_files(legacy) == [("database", SENTINEL)], (
    "Legacy Android backup must remain deny-by-default; review before allowing app data"
)

modern = ET.parse(ANDROID / "res/xml/data_extraction_rules.xml").getroot()
assert modern.tag == "data-extraction-rules", "Unexpected modern backup rules root"
for mode in ("cloud-backup", "device-transfer"):
    section = modern.find(mode)
    assert section is not None, f"Missing {mode} rules"
    assert allowed_files(section) == [("database", SENTINEL)], (
        f"{mode} must not export settings, profiles, or configuration files"
    )

assert not (ANDROID / "databases" / SENTINEL).exists(), "Backup sentinel must remain absent"
print("Android backup policy verified: no application databases are allowlisted.")

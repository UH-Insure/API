#!/usr/bin/env python3
"""
flattener.py

Usage:
    # Flatten a specific subfolder into all_files.txt at the repo root
    python flattener.py file/file
"""

import os
import argparse
from pathlib import Path

def write_all_files_here(
    out_name: str = "all_files.txt",
    include_hidden: bool = False,
    target_subpath: str | None = None
) -> None:
    # The "root" is where this script sits, and where the output will be written.
    root_dir = Path(__file__).resolve().parent
    out_path = root_dir / out_name

    # Determine which directory to scan:
    # - If target_subpath is provided, scan that subdir (relative to root)
    # - Otherwise, scan the root itself (legacy behavior)
    scan_dir = root_dir
    if target_subpath:
        # Normalize and ensure the scan_dir is inside root_dir
        candidate = (root_dir / target_subpath).resolve()
        try:
            candidate.relative_to(root_dir)
        except ValueError:
            raise ValueError(
                f"Refusing to scan outside the root: {candidate} is not within {root_dir}"
            )
        if not candidate.exists() or not candidate.is_dir():
            raise FileNotFoundError(f"Target directory does not exist: {candidate}")
        scan_dir = candidate

    with out_path.open("w", encoding="utf-8", errors="replace") as out_f:
        for dirpath, dirnames, filenames in os.walk(scan_dir):
            # Skip hidden directories/files unless requested
            if not include_hidden:
                dirnames[:] = [d for d in dirnames if not d.startswith('.')]
                filenames = [f for f in filenames if not f.startswith('.')]

            for fname in filenames:
                file_path = Path(dirpath) / fname
                # Skip the output file itself (in case scan_dir == root)
                if file_path.resolve() == out_path.resolve():
                    continue

                # Write the file path (relative to root for readability)
                try:
                    rel_from_root = file_path.resolve().relative_to(root_dir)
                    out_f.write(f"{root_dir / rel_from_root}\n")
                except Exception:
                    # Fallback to absolute if relative fails for any reason
                    out_f.write(f"{file_path.resolve()}\n")

                # Try to read as text
                try:
                    with file_path.open("r", encoding="utf-8", errors="replace") as f:
                        for line in f:
                            out_f.write(line)
                except Exception as e:
                    out_f.write(f"<ERROR: cannot read file: {e}>\n")

                out_f.write("\n")  # blank line between entries

    print(f"Wrote flattened file list to: {out_path}")
    if scan_dir != root_dir:
        print(f"Scanned directory: {scan_dir} (relative to root: {scan_dir.relative_to(root_dir)})")


def main():
    parser = argparse.ArgumentParser(description="Flatten files/contents into a single text file at the repo root.")
    parser.add_argument(
        "target_subpath",
        nargs="?",
        default=None,
        help="Relative path from the script's directory to the folder to scan (e.g. Backend/api-gateway-go). "
             "If omitted, the root directory (where this script lives) is scanned."
    )
    parser.add_argument(
        "--out-name",
        default="all_files.txt",
        help="Name of the output file written at the repo root (default: all_files.txt)."
    )
    parser.add_argument(
        "--include-hidden",
        action="store_true",
        help="Include hidden files and directories (those starting with '.')."
    )
    args = parser.parse_args()

    write_all_files_here(
        out_name=args.out_name,
        include_hidden=args.include_hidden,
        target_subpath=args.target_subpath
    )

if __name__ == "__main__":
    main()
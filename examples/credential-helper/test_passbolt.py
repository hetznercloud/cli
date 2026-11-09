"""Run with: python3 -m unittest discover -s examples/credential-helper."""

import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest


@unittest.skipIf(os.name == "nt", "The fake executable uses a POSIX shebang")
class PassboltHelperTest(unittest.TestCase):
    def test_lookup(self):
        for mode in ("success", "failure", "missing", "malformed", "multiline"):
            with self.subTest(mode=mode), tempfile.TemporaryDirectory() as directory:
                fake = Path(directory) / "passbolt"
                fake.write_text(
                    f"#!{sys.executable}\n"
                    "import json, os, sys\n"
                    "assert sys.argv[1:] == ['get', 'resource', '--id', "
                    "'01234567-89ab-cdef-0123-456789abcdef', '--json']\n"
                    "mode = os.environ['FAKE_PASSBOLT_MODE']\n"
                    "if mode == 'failure':\n"
                    "    print('private failure details', file=sys.stderr)\n"
                    "    sys.exit(1)\n"
                    "if mode == 'malformed':\n"
                    "    print('invalid JSON')\n"
                    "else:\n"
                    "    print(json.dumps({} if mode == 'missing' else "
                    "{'password': 'token\\ninvalid' if mode == 'multiline' else 'test-token'}))\n"
                )
                fake.chmod(0o700)
                result = subprocess.run(
                    [sys.executable, str(Path(__file__).with_name("passbolt.py")), "get"],
                    input="protocol=https\nhost=api.hetzner.cloud\n"
                    "path=01234567-89ab-cdef-0123-456789abcdef\n\n",
                    text=True,
                    capture_output=True,
                    env={**os.environ, "PATH": directory + os.pathsep + os.environ["PATH"],
                         "FAKE_PASSBOLT_MODE": mode},
                    check=False,
                )
                if mode == "success":
                    self.assertEqual(result.returncode, 0, result.stderr)
                    self.assertEqual(result.stdout, "password=test-token\n\n")
                else:
                    self.assertNotEqual(result.returncode, 0)
                    self.assertEqual(result.stdout, "")
                    self.assertNotIn("private failure details", result.stderr)


if __name__ == "__main__":
    unittest.main()

"""Sphinx configuration for the amortization-schedule documentation."""

project = "amortization-schedule"
author = "Michael Lip"
copyright = "2026, Michael Lip"
release = "0.1"

extensions = []
templates_path = []
exclude_patterns = ["_build"]

html_theme = "alabaster"
html_static_path = []
html_title = "amortization-schedule"
html_theme_options = {
    "description": "Fixed-rate loan payment and amortization schedule math for Go.",
    "github_user": "theluckystrike",
    "github_repo": "amortization-schedule",
    "github_button": False,
    "fixed_sidebar": True,
}

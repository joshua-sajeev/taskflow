## Understanding `templates/`

The `templates/` directory contains HTML templates for Taskflow's frontend.

### `index.tmpl`
- Provides a simple user interface for submitting jobs.
- Contains a form with an input field for entering a task and a submit button.
- Uses JavaScript to send job creation requests to the backend via the `/jobs` API.
- Displays success or error messages based on the server response.

This template allows users to interact with Taskflow's job management system through a web-based interface.


# ASCII Art Web – Implementation Exercise

## 1. Project Overview

This project extends the existing ASCII Art CLI application into a web application while preserving the current rendering engine. The banner loading, input normalization, validation, and ASCII rendering logic implemented in the previous project should be reused wherever possible. The primary objective is to replace the command-line interface with an HTTP server and HTML frontend while maintaining the same rendering behavior.

This is an extension of the existing codebase, **not a rewrite**. The rendering pipeline should remain the single source of truth for ASCII art generation.

---

# 2. Objectives

Implement a Go HTTP server that provides a graphical interface for generating ASCII art.

The application must:

- Serve an HTML webpage.
- Accept user input through an HTML form.
- Allow selection of one of the supported banners:
  - standard
  - shadow
  - thinkertoy
- Generate ASCII art using the existing renderer.
- Display the rendered ASCII art in the browser.
- Return appropriate HTTP status codes.
- Use only the Go standard library.

---

# 3. Existing Project Summary

The previous project already provides the core functionality required for rendering ASCII art.

Existing features include:

- Banner loading and validation.
- ASCII rendering.
- Escaped newline normalization.
- Support for:
  - standard
  - shadow
  - thinkertoy
- Input validation.
- Error handling for invalid banners and unsupported characters.

The CLI implementation should be treated as the reference implementation for rendering behavior.

Only the interface changes—from terminal interaction to HTTP requests.

The rendering logic itself should remain unchanged.

---

# 4. Reusable Components

Reuse the existing implementation wherever possible.

Recommended components to keep:

- Banner loader
- Banner validation
- ASCII renderer
- Input normalization
- Existing error handling related to rendering

Avoid duplicating rendering logic.

The HTTP handlers should simply collect user input, invoke the existing rendering pipeline, and return the generated result.

---

# 5. Suggested Project Architecture

A clean project structure might resemble:

```
.
├── main.go
├── handlers.go
├── templates/
│   └── index.html
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
├── internal/
│   ├── banner/
│   ├── renderer/
│   └── ...
└── README.md
```

Suggested responsibility breakdown:

### main.go

- Start the HTTP server.
- Register routes.
- Handle server startup errors.

### handlers.go

Implement all HTTP handlers.

Responsibilities include:

- Parsing requests.
- Validating form input.
- Calling the renderer.
- Executing templates.
- Returning appropriate HTTP status codes.

### templates/

Contains all HTML templates used by the application.

The assignment requires HTML templates to reside inside the `templates` directory.

---

# 6. Request Flow

```
Browser
    │
    ▼
GET /
    │
    ▼
Load HTML template
    │
    ▼
Display form
    │
    ▼
User submits form
    │
    ▼
POST /ascii-art
    │
    ▼
Read form values
    │
    ▼
Validate input
    │
    ▼
Load selected banner
    │
    ▼
Render ASCII art
    │
    ▼
Execute template
    │
    ▼
Return HTML response
```

---

# 7. HTTP API Specification

## GET /

Purpose:

Serve the main HTML page.

Expected behavior:

- Load the HTML template.
- Display:
  - text input
  - banner selector
  - submit button
  - rendered result (if applicable)

Successful response:

```
200 OK
```

Possible errors:

```
404 Not Found
```

If the template cannot be located.

```
500 Internal Server Error
```

For unexpected failures while executing templates.

---

## POST /ascii-art

Purpose:

Generate ASCII art from submitted form data.

Expected form values:

```
text
banner
```

Workflow:

1. Parse form.
2. Validate request.
3. Validate banner.
4. Render ASCII art.
5. Execute template with result.

Successful response:

```
200 OK
```

Possible errors:

```
400 Bad Request
```

Invalid or malformed request.

```
404 Not Found
```

Missing banner file.

```
500 Internal Server Error
```

Unexpected rendering or template execution failures.

---

# 8. HTML Requirements

The main page must contain:

- A text input field (or textarea).
- Banner selector.
  - Radio buttons
  - Select element
  - Any equivalent method
- Submit button.
- HTML form using POST.
- Display area for generated ASCII art.

Recommended HTML structure:

```
<form method="POST" action="/ascii-art">

    Text input

    Banner selection

    Submit button

</form>

<pre>

Generated ASCII Art

</pre>
```

Using a `<pre>` element is recommended to preserve spacing and formatting.

---

# 9. Template Data

The HTML template should receive data from the server using Go templates.

Example fields:

```
InputText

SelectedBanner

ASCIIArt

Error
```

Handlers should populate these values before executing the template.

---

# 10. Implementation Plan

## Phase 1

Set up the HTTP server.

Tasks:

- Create main.go.
- Register routes.
- Start listening on a port.

---

## Phase 2

Implement GET /

Tasks:

- Load HTML template.
- Render home page.
- Handle missing template errors.

---

## Phase 3

Create the HTML interface.

Tasks:

- Text input.
- Banner selector.
- Submit button.
- Display area.

---

## Phase 4

Implement POST /ascii-art

Tasks:

- Parse form.
- Validate request.
- Validate banner.
- Invoke renderer.
- Return HTML response.

---

## Phase 5

Integrate existing renderer.

Tasks:

- Reuse current banner loader.
- Reuse normalization.
- Reuse renderer.
- Avoid duplicating rendering code.

---

## Phase 6

Implement error handling.

Return correct HTTP status codes for:

- Bad requests.
- Missing resources.
- Internal failures.

---

## Phase 7

Create README.md

Include:

- Description
- Authors
- Usage
- Implementation details

---

# 11. HTTP Status Codes

The server must return appropriate status codes.

| Situation | Status |
|-----------|--------|
| Successful request | 200 OK |
| Invalid request | 400 Bad Request |
| Missing template | 404 Not Found |
| Missing banner | 404 Not Found |
| Unexpected error | 500 Internal Server Error |

The application should never panic during normal operation.

Errors should be handled gracefully with meaningful HTTP responses.

---

# 12. Testing Checklist

Verify the following:

- GET / returns HTML.
- Home page loads successfully.
- Form submits correctly.
- Standard banner renders correctly.
- Shadow banner renders correctly.
- Thinkertoy banner renders correctly.
- Escaped newlines behave as expected.
- Missing banner returns 404.
- Missing template returns 404.
- Invalid POST request returns 400.
- Unexpected failures return 500.
- ASCII art formatting is preserved.
- Browser output matches CLI rendering.
- README is present.

---

# 13. Constraints

Implementation must satisfy the following:

- Use only the Go standard library.
- HTML templates must be located inside the `templates` directory.
- Reuse the existing rendering pipeline.
- Do not duplicate renderer functionality.
- Do not modify banner files unless explicitly required.
- Follow Go best practices.
- Keep rendering logic separate from HTTP handling.

---

# 14. Deliverables

At completion, the project should include at minimum:

```
main.go

handlers.go

templates/
    index.html

README.md

banners/

internal/
```

The application should successfully:

- Serve an HTML interface.
- Accept user input.
- Render ASCII art using any supported banner.
- Display the generated ASCII art in the browser.
- Return appropriate HTTP status codes.
- Preserve the behavior of the existing ASCII rendering implementation while replacing only the user interface layer.
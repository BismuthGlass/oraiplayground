# ORAI Playground

A lightweight OpenRouter client using the orai library and HTMX.

## Structure

The project is structured into multiple packages, with the main package being responsible for starting the server and initializing everything.

### Models

Models are the basic data blocks that the application uses. They are mostly innert and are used by other entities to communicate consistently with the rest of the system.

### Templates

Templates in this project are stored as `.html` files written in Go's `html/template`. The Template package provides well defined structures that are to be passed to each template, and provides helper functions to convert various other models into template data. This package also exposes a service which simplifies the lookup of templates.

### Services

Services expose different functionality with the system, but do not expose any external APIs. They are used to do most of the internal work, such as updating or querying data, handling timers or request pools, and so on. Services can make use of other services, and they present the core functionality of the system.

### Controllers

Controllers are responsible for exposing APIs to the frontend, and handling high level client requests. Most of their job is delegated to Controllers.

### Components

Components in the backend are the controllers for frontend components. Since the webapp uses HTMX, most of the UI state is reflected through Controllers, which has the benefit of also being stored for future sessions, without the need for cookies or any other kind of local storage. This greatly simplifies the application at no cost, since it is mostly meant to be hosted locally, and allows us to use Go "on the frontend" as well.

### Config

Handles configuration globals. Nothing special about it. Exposes configurations as global variables, some having default values.

### Utils

Miscelaneous utilities and helper functions used by all other packages.

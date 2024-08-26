# go_image_upload
upload images

```
repo/
│
├── cmd/
│   └── server/
│       └── main.go            # Entry point of the application
│
├── internal/
│   ├── handlers/
│   │   └── handlers.go        # Registering routes
│   │
│   ├── services/
│   │   └── file_service.go    # Business logic for file operations
│   │
│   ├── storage/
│   │   └── file_storage.go    # File system operations (saving, deleting, etc.)
│   │
│   └── config/
│       └── config.go          # Configuration settings (e.g., constants like uploadDir)
│
├── pkg/
│   └── utils/
│       └── utils.go           # Utility functions (e.g., UUID generation)
│
└── go.mod                     # Go module file
```
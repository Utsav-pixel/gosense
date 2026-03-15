# 🚀 GoSense v0.2.0 - Public API Release

## 🎯 **Major Release - Breaking Change**

This release fixes the critical issue where external users could not import the GoSense library due to Go's internal package restrictions.

## ✨ **What's New**

### **🔧 Public API**
- **NEW**: Public API package (`github.com/Utsav-pixel/gosense`) that re-exports all engine functionality
- **FIXED**: External users can now import and use the library with `go get github.com/Utsav-pixel/gosense`
- **UPDATED**: All examples now use the public API

### **📚 Documentation Overhaul**
- **NEW**: Comprehensive blog post explaining the library's purpose and usage
- **UPDATED**: All documentation files reference the public API
- **IMPROVED**: Installation instructions and usage examples
- **ADDED**: Quick start templates and best practices

### **🛠️ Installation**
```bash
# Now works for external users!
go get github.com/Utsav-pixel/gosense
```

### **📖 Usage Example**
```go
import "github.com/Utsav-pixel/gosense"

// Create seeder
seeder := gosense.NewTimeSeeder(1.0, 0.1, 20.0)

// Create sensor function
sensorFunc := gosense.NewFunction(func(input float64, timestamp time.Time) YourData {
    return YourData{Value: input * 100}
})

// Create engine
engine := gosense.NewEngine(config, seeder, sensorFunc, publisher)
```

## 🔄 **Breaking Changes**

- **Import Path**: Changed from `github.com/Utsav-pixel/gosense/internal/engine` to `github.com/Utsav-pixel/gosense`
- **API Prefix**: Changed from `engine.` to `gosense.`
- **Package Structure**: Internal packages are no longer accessible (as intended)

## 📋 **Migration Guide**

### **Before (v0.1.0)**
```go
import "github.com/Utsav-pixel/gosense/internal/engine"

seeder := engine.NewTimeSeeder(1.0, 0.1, 20.0)
```

### **After (v0.2.0)**
```go
import "github.com/Utsav-pixel/gosense"

seeder := gosense.NewTimeSeeder(1.0, 0.1, 20.0)
```

## 🎉 **Impact**

- ✅ **External users can now use the library**
- ✅ **Clean public API with full functionality**
- ✅ **Comprehensive documentation**
- ✅ **Proper Go package structure**
- ✅ **All examples working**

## 🔗 **Links**

- [GitHub Repository](https://github.com/Utsav-pixel/gosense)
- [Documentation](https://github.com/Utsav-pixel/gosense/tree/main/docs)
- [Examples](https://github.com/Utsav-pixel/gosense/tree/main/examples)

---

**Note**: This is a breaking change but enables the core functionality of the library. All existing functionality is preserved, just through the public API instead of internal packages.

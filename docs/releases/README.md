# 🚀 GoSense Release Notes

This directory contains detailed release notes for all GoSense versions.

## 📋 Available Releases

### **[v0.2.2](./v0.2.2.md)** - Documentation Organization
- ✅ Organized release notes in dedicated folder
- ✅ Enhanced README with latest release section
- ✅ Better documentation navigation and structure
- ✅ Professional documentation presentation

### **[v0.2.1](./v0.2.1.md)** - Complete Public API
- ✅ Added all publishers to public API (HTTP, Kafka, gRPC)
- ✅ Complete documentation updates
- ✅ Production-ready library

### **[v0.2.0](./v0.2.0.md)** - Public API Release
- ✅ Fixed critical import issue
- ✅ Added public API package
- ✅ Breaking change from internal packages

---

## 🔄 Upgrade Guide

### From v0.1.0 to v0.2.2
```bash
# Old (broken)
import "github.com/Utsav-pixel/gosense/internal/engine"

# New (working)
import "github.com/Utsav-pixel/gosense"
```

### From v0.2.0/v0.2.1 to v0.2.2
No breaking changes - documentation improvements only:
```bash
go get github.com/Utsav-pixel/gosense@v0.2.2
```

---

## 📈 Release History

| Version | Date | Changes |
|---------|------|---------|
| v0.2.2 | 2025-03-15 | Documentation organization |
| v0.2.1 | 2025-03-15 | Complete public API with publishers |
| v0.2.0 | 2025-03-15 | Fixed public API access |
| v0.1.0 | - | Initial release (broken for external users) |

---

**🎯 Latest Version**: [v0.2.2](./v0.2.2.md)  
**📦 Install**: `go get github.com/Utsav-pixel/gosense@v0.2.2`

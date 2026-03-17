# 🚀 GoSense Release Notes

This directory contains detailed release notes for all GoSense versions.

## 📋 Available Releases

### **[v0.2.5](./v0.2.5.md)** - Blog Migration to Medium
- ✅ Moved blog from repository to Medium platform
- ✅ Reduced repository size by removing large blog file
- ✅ Updated all references to point to Medium article
- ✅ Better reading experience and accessibility

### **[v0.2.4](./v0.2.4.md)** - Real Performance Benchmarks
- ✅ Replaced made-up performance numbers with real benchmarks
- ✅ Added complete benchmarking suite for verification
- ✅ Updated all documentation with actual performance metrics
- ✅ Added machine specifications and performance tips

### **[v0.2.3](./v0.2.3.md)** - Console Publisher Added
- ✅ Added console publisher to public API
- ✅ Complete publisher ecosystem (HTTP, Kafka, gRPC, Console)
- ✅ Enhanced development and testing experience
- ✅ Visual formatted output for debugging

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

### From v0.1.0 to v0.2.5
```bash
# Old (broken)
import "github.com/Utsav-pixel/gosense/internal/engine"

# New (working)
import "github.com/Utsav-pixel/gosense"
```

### From v0.2.0/v0.2.1/v0.2.2/v0.2.3/v0.2.4 to v0.2.5
No breaking changes - documentation improvements only:
```bash
go get github.com/Utsav-pixel/gosense@v0.2.5
```

---

## 📈 Release History

| Version | Date | Changes |
|---------|------|---------|
| v0.2.5 | 2025-03-18 | Blog migration to Medium |
| v0.2.4 | 2025-03-16 | Real performance benchmarks |
| v0.2.3 | 2025-03-15 | Console publisher added |
| v0.2.2 | 2025-03-15 | Documentation organization |
| v0.2.1 | 2025-03-15 | Complete public API with publishers |
| v0.2.0 | 2025-03-15 | Fixed public API access |
| v0.1.0 | - | Initial release (broken for external users) |

---

**🎯 Latest Version**: [v0.2.5](./v0.2.5.md)  
**📦 Install**: `go get github.com/Utsav-pixel/gosense@v0.2.5`

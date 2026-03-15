# 🚀 GoSense v0.2.1 - Complete Public API

## 🎯 **Enhancement Release**

This release completes the public API by adding all built-in publishers that were previously missing.

## ✨ **What's New**

### **📡 Complete Publisher Support**
- **NEW**: `gosense.NewGenericHTTPPublisher[T]()` - HTTP/REST API publishing
- **NEW**: `gosense.NewGenericKafkaPublisher[T]()` - Apache Kafka publishing
- **NEW**: `gosense.NewGenericGRPCPublisher[T]()` - gRPC streaming

### **📚 Updated Documentation**
- **UPDATED**: All examples now show publisher usage
- **IMPROVED**: README includes complete publisher examples
- **ENHANCED**: Advanced guide with publisher configurations
- **FIXED**: Blog post uses public API for all examples

## 📖 **Complete Usage Example**

```go
import "github.com/Utsav-pixel/gosense"

// Create seeder
seeder := gosense.NewTimeSeeder(1.0, 0.1, 20.0)

// Create sensor function
sensorFunc := gosense.NewFunction(func(input float64, timestamp time.Time) YourData {
    return YourData{Value: input * 100}
})

// Choose your publisher
httpPublisher := gosense.NewGenericHTTPPublisher[YourData]("https://api.example.com/data")
// OR
kafkaPublisher := gosense.NewGenericKafkaPublisher[YourData]([]string{"localhost:9092"}, "topic")
// OR  
grpcPublisher, err := gosense.NewGenericGRPCPublisher[YourData]("localhost:50051")

// Create engine
engine := gosense.NewEngine(config, seeder, sensorFunc, httpPublisher)
```

## 🔧 **Publisher Details**

### **HTTP Publisher**
- Simple REST API integration
- JSON payload format
- Configurable timeout (5s default)
- Batch and single message support

### **Kafka Publisher**
- High-throughput streaming
- Automatic message keying (uses sensor ID)
- Configurable brokers and topics
- Built-in connection pooling

### **gRPC Publisher**
- Real-time streaming
- Binary protocol efficiency
- Configurable server address
- Secure connection support

## 🎉 **Impact**

- ✅ **Complete public API** - All functionality now accessible
- ✅ **Production-ready publishers** - HTTP, Kafka, gRPC all available
- ✅ **Comprehensive examples** - Clear usage patterns documented
- ✅ **Zero internal dependencies** - Clean external interface

## 🔄 **Changes from v0.2.0**

- Added publisher constructors to public API
- Updated all documentation examples
- No breaking changes - fully backward compatible

## 🔗 **Links**

- [GitHub Repository](https://github.com/Utsav-pixel/gosense)
- [Documentation](https://github.com/Utsav-pixel/gosense/tree/main/docs)
- [Examples](https://github.com/Utsav-pixel/gosense/tree/main/examples)

---

**Note**: This release completes the public API. Users now have access to all core functionality including seeders, functions, engines, and publishers through a clean, documented interface.

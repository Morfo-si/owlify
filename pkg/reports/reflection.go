package reports

import (
    "fmt"
    "reflect"
    "strings"
    "time"
)

// capitalizeFirst capitalizes the first letter of a string
func CapitalizeFirst(s string) string {
    if s == "" {
        return ""
    }
    return strings.ToUpper(s)
}

// GetFlattenedHeaders returns a flattened list of headers from a struct
func GetFlattenedHeaders(data interface{}) []string {
    var t reflect.Type
    
    // Handle both reflect.Type and interface{} inputs
    if rt, ok := data.(reflect.Type); ok {
        t = rt
    } else {
        t = reflect.TypeOf(data)
        if t.Kind() == reflect.Slice {
            t = t.Elem()
        }
    }
    
    var headers []string
    extractHeaders(t, "", &headers)
    return headers
}

// extractHeaders recursively extracts header names from struct fields
func extractHeaders(t reflect.Type, prefix string, headers *[]string) {
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        if prefix != "" {
            *headers = append(*headers, prefix)
        }
        return
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        
        // Skip unexported fields
        if field.PkgPath != "" {
            continue
        }
        
        // Get field name from JSON tag if available
        fieldName := field.Name
        jsonTag := field.Tag.Get("json")
        if jsonTag != "" && jsonTag != "-" {
            parts := strings.Split(jsonTag, ",")
            if parts[0] != "" {
                fieldName = parts[0]
            }
        }
        
        // Always use lowercase field names
        fieldName = strings.ToLower(fieldName)
        
        // Build the full path for this field
        fullPath := fieldName
        if prefix != "" {
            fullPath = prefix + "." + fieldName
        }
        
        // Special handling for Epic and Feature types
        if field.Type.Kind() == reflect.Ptr && 
           (field.Type.Elem().Name() == "Epic" || field.Type.Elem().Name() == "Feature") {
            *headers = append(*headers, fullPath+".key", fullPath+".summary")
            continue
        }
        
        // Handle nested structs recursively
        if field.Type.Kind() == reflect.Struct {
            // Check if this is a special type that shouldn't be flattened
            if field.Type.Name() == "Time" {
                *headers = append(*headers, fullPath)
            } else {
                extractHeaders(field.Type, fullPath, headers)
            }
        } else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
            extractHeaders(field.Type.Elem(), fullPath, headers)
        } else {
            *headers = append(*headers, fullPath)
        }
    }
}

// GetFlattenedHeadersRecursive recursively gets headers from a struct type
func GetFlattenedHeadersRecursive(t reflect.Type, prefix string, headers *[]string, skipUnexported bool) {
    if t == nil {
        return
    }
    
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        if prefix != "" {
            *headers = append(*headers, prefix)
        }
        return
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        
        // Skip unexported fields if requested
        if skipUnexported && field.PkgPath != "" {
            continue
        }
        
        jsonTag := field.Tag.Get("json")
        if jsonTag == "-" {
            continue
        }
        
        fieldName := field.Name
        if jsonTag != "" {
            parts := strings.Split(jsonTag, ",")
            if parts[0] != "" {
                fieldName = parts[0]
            }
        }
        
        // Skip "fields" struct level
        if strings.ToLower(fieldName) == "fields" {
            GetFlattenedHeadersRecursive(field.Type, prefix, headers, skipUnexported)
            continue
        }
        
        // Special handling for Epic and Feature types
        if field.Type.Kind() == reflect.Ptr && 
           (field.Type.Elem().Name() == "Epic" || field.Type.Elem().Name() == "Feature") {
            if prefix != "" {
                fieldName = prefix + "." + fieldName
            }
            *headers = append(*headers, fieldName+".Key", fieldName+".Summary")
            continue
        }
        
        newPrefix := fieldName
        if prefix != "" {
            newPrefix = prefix + "." + fieldName
        }
        
        if field.Type.Kind() == reflect.Struct {
            GetFlattenedHeadersRecursive(field.Type, newPrefix, headers, skipUnexported)
        } else {
            *headers = append(*headers, newPrefix)
        }
    }
}

// getFlattenedValues returns a slice of values for all fields including nested structs
func getFlattenedValues(v reflect.Value) []string {
    var values []string
    t := v.Type()
    
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        structField := t.Field(i)
        
        // Skip unexported fields
        if !structField.IsExported() {
            continue
        }
        
        jsonTag := structField.Tag.Get("json")
        if jsonTag == "-" {
            continue
        }

        fieldName := structField.Name
        if jsonTag != "" {
            parts := strings.Split(jsonTag, ",")
            if parts[0] != "" {
                fieldName = parts[0]
            }
        }

        // Skip "fields" struct level
        if fieldName == "fields" {
            nestedValues := getFlattenedValues(field)
            values = append(values, nestedValues...)
            continue
        }

        // Handle special pointer struct fields
        if IsSpecialPointerField(fieldName, field) {
            values = append(values, getSpecialFieldValues(field)...)
            continue
        }

        if field.Kind() == reflect.Struct {
            // For embedded fields, include their values directly
            if structField.Anonymous {
                nestedValues := getFlattenedValues(field)
                values = append(values, nestedValues...)
            } else {
                // Recursively get values for nested struct
                nestedValues := getFlattenedValues(field)
                values = append(values, nestedValues...)
            }
        } else {
            // Handle nil pointers
            if field.Kind() == reflect.Ptr && field.IsNil() {
                values = append(values, "")
                continue
            }
            
            // Get value (dereference pointer if needed)
            var val interface{}
            if field.Kind() == reflect.Ptr {
                val = field.Elem().Interface()
            } else {
                val = field.Interface()
            }
            
            // Special formatting for time.Time values
            switch v := val.(type) {
            case time.Time:
                // Format time as YYYY-MM-DD
                values = append(values, v.Format("2006-01-02"))
            default:
                values = append(values, fmt.Sprintf("%v", val))
            }
        }
    }
    
    return values
}

// Helper functions for special field handling

// isSpecialStructField checks if a field is a special case that needs custom handling
func IsSpecialStructField(_ string, field reflect.StructField) bool {
    // Add special case handling for time.Time
    return field.Type.String() == "time.Time"
}

// IsSpecialPointerField checks if a field is a special pointer type that needs special handling
func IsSpecialPointerField(fieldName string, field reflect.Value) bool {
    // Check if it's a pointer and not nil
    if field.Kind() == reflect.Ptr && !field.IsNil() {
        // Check if it's an Epic or Feature type
        typeName := field.Elem().Type().Name()
        return typeName == "Epic" || typeName == "Feature"
    }
    return false
}

// getSpecialFieldHeaders returns headers for special field types
func GetSpecialFieldHeaders(fieldName string, field reflect.StructField, prefix string) []string {
    if field.Type.String() == "time.Time" {
        if prefix != "" {
            return []string{prefix + fieldName}
        }
        return []string{fieldName}
    }
    return []string{fieldName}
}

// getSpecialFieldValues extracts values from special fields like Epic and Feature
func getSpecialFieldValues(field reflect.Value) []string {
    if field.Kind() != reflect.Ptr || field.IsNil() {
        return []string{""}
    }
    
    elem := field.Elem()
    typeName := elem.Type().Name()
    
    switch typeName {
    case "Epic":
        // Extract Key and Summary from Epic
        key := elem.FieldByName("Key").String()
        summary := elem.FieldByName("Summary").String()
        return []string{key, summary}
    case "Feature":
        // Extract Key and Summary from Feature
        key := elem.FieldByName("Key").String()
        summary := elem.FieldByName("Summary").String()
        return []string{key, summary}
    default:
        return []string{fmt.Sprintf("%v", field.Interface())}
    }
}

// buildNestedPrefix creates a prefix for nested struct fields
func BuildNestedPrefix(prefix, fieldName string) string {
    if prefix != "" {
        return prefix + fieldName + " . "
    }
    return fieldName + " . "
}

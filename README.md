# Locale detection for caddy

This plugin matches preferred user locales to available locales configured on server side.  
User preferences can be detected from the `Accept-Language` header or a cookie.  
It sets the header defined by `header` to the detected value (default header name `Detected-Locale`).  

## Configuration schema

```
locale <availableLocales...> {
  available <availableLocales...>
  methods <methods...>
  cookie <cookie name>
  header <header name>
}
```

A `method` can be currently `cookie` or `header`.  
If `cookie` is added, `cookie name` defines from which cookie the locale is read.  
The `header` method extracts the locales from the `Accept-Language` header.

The first `availableLocale` is also the default locale, which is picked if none of the locales from the detection methods is in `availableLocales`.

The defaults are: `methods = [header]`,  `cookie name = lang`, `header name = Detected-Locale`.

## Examples

### Minimal example

```
locale en de 
```

### Override / normalize `Accept-Language` header

```
locale en_US en_GB de {
  header Accept-Language
}
```

### Read locale from cookie

```
locale de en es fr {
  methods cookie
}
```

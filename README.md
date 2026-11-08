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
  region_header <header name>
}
```

A `method` can be currently `cookie` or `header`.  
If `cookie` is added, `cookie name` defines from which cookie the locale is read.  
The `header` method extracts the locales from the `Accept-Language` header.

The first `availableLocale` is also the default locale, which is picked if none of the locales from the detection methods is in `availableLocales`.

The defaults are: `methods = [header]`, `cookie name = lang`, `header name = Detected-Locale`.

`region_header` (off by default) names a header that receives the region the client most likely lives in, read from
the original `Accept-Language` (before `header` may replace it):

- the region of the first accepted language that names a country: `fr-CA,fr;q=0.9` gives `CA`;
- else the likely region of the preferred language (CLDR likely subtags), marked as a guess: `de` gives `DE; guessed`,
  `pt` gives `BR; guessed`;
- else nothing (`*`, no header).

The header is always written or removed, so a client cannot send its own. The region is also available as
`{vars.detected-region}`. It is a hint from the browser's settings, not a location.

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

### Detect the client's region too

```
locale en fr de {
  header Accept-Language
  region_header Detected-Region
}
```

`Accept-Language: fr-CA,fr;q=0.9` reaches the backend as `Accept-Language: fr` and `Detected-Region: CA`.

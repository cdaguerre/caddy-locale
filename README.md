# Locale detection for caddy

## Configuration schema

    locale <availableLocales...> {
      available <availableLocales...>
      methods <methods...>
      cookie <cookie name>
    }

A `method` can be currently `cookie` or `header`. If `cookie` is added, `cookie name` defines from which cookie the
locale is read. The `header` method extracts the locales from the `Accept-Language` header. The first `availableLocale`
is also the default locale, which is picked if none of the locales from the detection methods is in `availableLocales`.

The defaults are: `methods = [header]`,  `cookie name = lang`.

## Example

    locale en de {
      detect cookie header
    }
# Changes

## 2.0.0

* Add support for `RETURNING` clause
* Quote column names
* Add MapSQL function which does not quote the argument
* Fix scanning into nil interface

By quoting column names, 2.0.0 may break programs using version 1. Queries
with `Map("COUNT(*)", &count)` will not work as expected due to quoting
the argument. Use the new MapSQL function in this case: `MapSQL("COUNT(*)", &count)`

## 1.0.1

* Surround where conditions in parenthesis

## 1.0.0

* Initial release

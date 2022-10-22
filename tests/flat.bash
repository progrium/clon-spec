#!/bin/bash

PROG=./clon

T_flatString() {
  result="$($PROG \
    name=John \
    hair=short \
    2>/dev/null \
  )"
  [[ "${result}" == '{"name":"John","hair":"short"}' ]]
}

T_flatStringBool() {
  result="$($PROG \
    name=John \
    live:=true \
    2>/dev/null \
  )"
  [[ "${result}" == '{"name":"John","live":true}' ]]
}

T_flatStringNumber() {
  result="$($PROG \
    name=John \
    age:=28 \
    2>/dev/null \
  )"
  [[ "${result}" == '{"name":"John","age":28}' ]]
}

T_flatRawArray() {
  result="$($PROG \
    name=John \
    languages:='["english","zulu"]' \
    2>/dev/null \
  )"
  [[ "${result}" == '{"name":"John","languages":["english","zulu"]}' ]]
}

T_flatRawObj() {
  result="$($PROG \
    name=John \
    wife:='{"name":"Jane","age":30}' \
    2>/dev/null \
  )"
  [[ "${result}" == '{"name":"John","wife":{"name":"Jane","age":30}}' ]]
}

T_flatAll() {
  result="$($PROG \
    name=John \
    age:=28 \
    lang:='["english","zulu"]' \
    married:=true \
    wife:='{"name":"Jane","age":30}' \
    2>/dev/null \
  )"
  [[ "${result}" == '{"name":"John","age":28,"lang":["english","zulu"],"married":true,"wife":{"name":"Jane","age":30}}' ]]
}


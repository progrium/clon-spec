#!/bin/bash

PROG=./clon

T_arrayFlatSingle() {
    result="$($PROG \
    name=John \
    lang[]=english \
  )"

  local expected='{"name":"John","lang":["english"]}'
  [[ "${result}" == "${expected}" ]] # || $T_fail "not matching"
}

T_arrayFlatMulti() {
    result="$($PROG \
    name=John \
    lang[]=english \
    lang[]=zulu \
    lang[]=german \
  )"

  local expected='{"name":"John","lang":["english","zulu","german"]}'
  [[ "${result}" == "${expected}" ]] # || $T_fail "not matching"
}

T_arrayNestedSingle() {
    result="$($PROG \
    name=John \
    wife[lang][]=english \
  )"

  local expected='{"name":"John","wife":{"lang":["english"]}}'
  [[ "${result}" == "${expected}" ]] # || $T_fail "not matching"
}

T_arrayNestedMulti() {
    result="$($PROG \
    name=John \
    lang[]=english \
    lang[]=zulu \
    wife[lang][]=french \
    wife[lang][]=spanish \
    wife[lang][]=greek \
  )"

  local expected='{"name":"John","lang":["english","zulu"],"wife":{"lang":["french","spanish","greek"]}}'
  [[ "${result}" == "${expected}" ]] # || $T_fail "not matching"
}

T_arrayMixingRawAndNestedMulti() {
    result="$($PROG \
    category=tools \
    search[type]=platforms \
    search[platforms]:='["Terminal", "Desktop"]' \
    search[platforms][]=Web \
    search[platforms][]=Mobile \
  )"

  local expected='{"category":"tools","search":{"type":"platforms","platforms":["Terminal","Desktop","Web","Mobile"]}}'
  [[ "${result}" == "${expected}" ]] # || $T_fail "not matching"
}




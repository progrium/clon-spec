#!/bin/bash

PROG=./clon

T_nestedString() {
  result="$($PROG \
    wife[name]=Jane \
  )"
  expected='{"wife":{"name":"Jane"}}'
  [[ "${result}" == "${expected}" ]]
}

T_nestedNumber() {
  result="$($PROG \
    wife[age]:=30 \
  )"
    expected='{"wife":{"age":30}}'
  [[ "${result}" == "${expected}" ]]
}

T_nestedStringAndNumber() {
  result="$($PROG \
    wife[name]=Jane \
    wife[age]:=30 \
  )"
  expected='{"wife":{"name":"Jane","age":30}}'
  [[ "${result}" == "${expected}" ]]
}

T_nestedStringNumberBool() {
    result="$($PROG \
    wife[name]=Jane \
    wife[age]:=30 \
    wife[female]:=true \
  )"
  expected='{"wife":{"name":"Jane","age":30,"female":true}}'
  [[ "${result}" == "${expected}" ]]
}

T_nestedMultiMixed() {
    result="$($PROG \
    wife[name]]=Jane \
    wife[age]:=30 \
    mother[name]=Tereza \
    wife[female]:=true \
    mother[age]:=70 \
    mother[female]:=true \
  )"
  expected='{"wife":{"name":"Jane","age":30,"female":true},"mother":{"name":"Tereza","age":70,"female":true}}'
  [[ "${result}" == "${expected}" ]]
}
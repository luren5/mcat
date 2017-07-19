pragma solidity ^0.4.13;

contract Foo {
  function bar(fixed[2] xy) {}
  function baz(uint32 x, bool y) returns (bool r) { r = x > 32 || y; }
  function sam(bytes name, bool z, uint[] data) {}
  function f(uint,uint32[],bytes10,bytes){}
}


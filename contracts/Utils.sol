pragma solidity ^0.4.23;


library Utils {

    function unpackSignatureAndSigns(bytes32 inputParam) internal pure returns(address addr, uint8 v1, uint8 v2) {
        addr = address(inputParam);
        v1 = uint8(inputParam[10]);
        v2 = uint8(inputParam[11]);
        return;
    }
}
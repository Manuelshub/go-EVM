package types

type Opcode byte

// Stop opcode
const (
	STOP Opcode = 0x00 // Stop and halts execution
)

// Arithmetic & comparison opcodes
const (
	ADD        Opcode = 0x01 // Addition operation
	MUL        Opcode = 0x02 // Multiplication operation
	SUB        Opcode = 0x03 // Subtraction operation
	DIV        Opcode = 0x04 // Integer division operation
	SDIV       Opcode = 0x05 // Signed integer division operation
	MOD        Opcode = 0x06 // Modulo remainder operation
	SMOD       Opcode = 0x07 // Signed modulo remainder operation
	ADDMOD     Opcode = 0x08 // Modulo addition operation
	MULMOD     Opcode = 0x09 // Modulo multiplication operation
	EXP        Opcode = 0x0A // Exponential operation
	SIGNEXTEND Opcode = 0x0B // Extend length of two's compliment signed integer
	LT         Opcode = 0x10 // Less-than comparison
	GT         Opcode = 0x11 // Greater-than comparison
	SLT        Opcode = 0x12 // Signed less-than comparison
	SGT        Opcode = 0x13 // Signed greater-than comparison
	EQ         Opcode = 0x14 // Equality comparison
	ISZERO     Opcode = 0x15 // Is-zero comparison
	AND        Opcode = 0x16 // Bitwise AND operation
	OR         Opcode = 0x17 // Bitwise OR operation
	XOR        Opcode = 0x18 // Bitwise XOR operation
	NOT        Opcode = 0x19 // Bitwise NOT operation
	BYTE       Opcode = 0x1A // Retrieve single byte from word
	SHL        Opcode = 0x1B // Left shift operation
	SHR        Opcode = 0x1C // Right shift operation
	SAR        Opcode = 0x1D // Arithmetic (signed) right shift operation
)

// Environmental opcodes
const (
	PC             Opcode = 0x58 // Get the value of the program counter prior to the increment
	GAS            Opcode = 0x5A // Get the amount of available gas
	KECCAK256      Opcode = 0x20 // Compute Keccak-256 hash
	ADDRESS        Opcode = 0x30 // Get address of currently executing account
	BALANCE        Opcode = 0x31 // Get balance of the given account
	ORIGIN         Opcode = 0x32 // Get execution origination address
	CALLER         Opcode = 0x33 // Get caller address
	CALLVALUE      Opcode = 0x34 // Get deposited value by the instruction/transaction responsible for this execution
	CALLDATALOAD   Opcode = 0x35 // Get input data of current environment
	CALLDATASIZE   Opcode = 0x36 // Get size of input data in current environment
	CALLDATACOPY   Opcode = 0x37 // Copy input data in current environment to memory
	CODESIZE       Opcode = 0x38 // Get size of code running in current environment
	CODECOPY       Opcode = 0x39 // Copy code running in current environment to memory
	GASPRICE       Opcode = 0x3A // Get price of gas in current environment
	EXTCODESIZE    Opcode = 0x3B // Get size of an account's code
	EXTCODECOPY    Opcode = 0x3C // Copy an account's code to memory
	RETURNDATASIZE Opcode = 0x3D // Get size of output data from the previous call from the current environment
	RETURNDATACOPY Opcode = 0x3E // Copy output data from the previous call to memory
	EXTCODEHASH    Opcode = 0x3F // Get the hash of account's code
	BLOCKHASH      Opcode = 0x40 // Get the hash of one of the 256 most recent complete blocks
	COINBASE       Opcode = 0x41 // Get the block's beneficiary address
	TIMESTAMP      Opcode = 0x42 // Get the block's timestamp
	NUMBER         Opcode = 0x43 // Get the block's number
	PREVRANDAO     Opcode = 0x44 // Get the block's difficulty
	GASLIMIT       Opcode = 0x45 // Get the block's gas limit
	CHAINID        Opcode = 0x46 // Get the chain's ID
	SELFBALANCE    Opcode = 0x47 // Get balance of currently executing account
	BASEFEE        Opcode = 0x48 // Get the block's base fee
	BLOBHASH       Opcode = 0x49 // Get versioned hashes
	BLOBBASEFEE    Opcode = 0x4A // Return the value of the blob base-fee of the current block
)

// Stack pop
const (
	POP = 0x50 // Remove item from stack
)

// Stack push opcodes
const (
	PUSH0  Opcode = 0x5F // Place 0 on the stack
	PUSH1  Opcode = 0x60 // Place 1 byte item on the stack
	PUSH2  Opcode = 0x61 // Place 2 byte item on the stack
	PUSH3  Opcode = 0x62 // Place 3 byte item on the stack
	PUSH4  Opcode = 0x63 // Place 4 byte item on the stack
	PUSH5  Opcode = 0x64 // Place 5 byte item on the stack
	PUSH6  Opcode = 0x65 // Place 6 byte item on the stack
	PUSH7  Opcode = 0x66 // Place 7 byte item on the stack
	PUSH8  Opcode = 0x67 // Place 8 byte item on the stack
	PUSH9  Opcode = 0x68 // Place 9 byte item on the stack
	PUSH10 Opcode = 0x69 // Place 10 byte item on the stack
	PUSH11 Opcode = 0x6A // Place 11 byte item on the stack
	PUSH12 Opcode = 0x6B // Place 12 byte item on the stack
	PUSH13 Opcode = 0x6C // Place 13 byte item on the stack
	PUSH14 Opcode = 0x6D // Place 14 byte item on the stack
	PUSH15 Opcode = 0x6E // Place 15 byte item on the stack
	PUSH16 Opcode = 0x6F // Place 16 byte item on the stack
	PUSH17 Opcode = 0x70 // Place 17 byte item on the stack
	PUSH18 Opcode = 0x71 // Place 18 byte item on the stack
	PUSH19 Opcode = 0x72 // Place 19 byte item on the stack
	PUSH20 Opcode = 0x73 // Place 20 byte item on the stack
	PUSH21 Opcode = 0x74 // Place 21 byte item on the stack
	PUSH22 Opcode = 0x75 // Place 22 byte item on the stack
	PUSH23 Opcode = 0x76 // Place 23 byte item on tHe stack
	PUSH24 Opcode = 0x77 // Place 24 byte item on the stack
	PUSH25 Opcode = 0x78 // Place 25 byte item on the stack
	PUSH26 Opcode = 0x79 // Place 26 byte item on the stack
	PUSH27 Opcode = 0x7A // Place 27 byte item on the stack
	PUSH28 Opcode = 0x7B // Place 28 byte item on the stack
	PUSH29 Opcode = 0x7C // Place 29 byte item on the stack
	PUSH30 Opcode = 0x7D // Place 30 byte item on the stack
	PUSH31 Opcode = 0x7E // Place 31 byte item on the stack
	PUSH32 Opcode = 0x7F // Place 32 byte (full word) item on the stack
)

// Stack Duplication opcodes
const (
	DUP1  Opcode = 0x80 // Duplicate 1st item on the stack
	DUP2  Opcode = 0x81 // Duplicate 2nd item on the stack
	DUP3  Opcode = 0x82 // Duplicate 3rd item on the stack
	DUP4  Opcode = 0x83 // Duplicate 4th item on the stack
	DUP5  Opcode = 0x84 // Duplicate 5th item on the stack
	DUP6  Opcode = 0x85 // Duplicate 6th item on the stack
	DUP7  Opcode = 0x86 // Duplicate 7th item on the stack
	DUP8  Opcode = 0x87 // Duplicate 8th item on the stack
	DUP9  Opcode = 0x88 // Duplicate 9th item on the stack
	DUP10 Opcode = 0x89 // Duplicate 10th item on the stack
	DUP11 Opcode = 0x8A // Duplicate 11th item on the stack
	DUP12 Opcode = 0x8B // Duplicate 12th item on the stack
	DUP13 Opcode = 0x8C // Duplicate 13th item on the stack
	DUP14 Opcode = 0x8D // Duplicate 14th item on the stack
	DUP15 Opcode = 0x8E // Duplicate 15th item on the stack
	DUP16 Opcode = 0x8F // Duplicate 16th item on the stack
)

// Stack Swap opcodes
const (
	SWAP1  Opcode = 0x90 // Exchange 1st and 2nd stack items
	SWAP2  Opcode = 0x91 // Exchange 1st and 3rd stack items
	SWAP3  Opcode = 0x92 // Exchange 1st and 4th stack items
	SWAP4  Opcode = 0x93 // Exchange 1st and 5th stack items
	SWAP5  Opcode = 0x94 // Exchange 1st and 6th stack items
	SWAP6  Opcode = 0x95 // Exchange 1st and 7th stack items
	SWAP7  Opcode = 0x96 // Exchange 1st and 8th stack items
	SWAP8  Opcode = 0x97 // Exchange 1st and 9th stack items
	SWAP9  Opcode = 0x98 // Exchange 1st and 10th stack items
	SWAP10 Opcode = 0x99 // Exchange 1st and 11th stack items
	SWAP11 Opcode = 0x9A // Exchange 1st and 12th stack items
	SWAP12 Opcode = 0x9B // Exchange 1st and 13th stack items
	SWAP13 Opcode = 0x9C // Exchange 1st and 14th stack items
	SWAP14 Opcode = 0x9D // Exchange 1st and 15th stack items
	SWAP15 Opcode = 0x9E // Exchange 1st and 16th stack items
	SWAP16 Opcode = 0x9F // Exchange 1st and 17th stack items
)

// Memory opcodes
const (
	MLOAD   Opcode = 0x51 // Load word from memory
	MSTORE  Opcode = 0x52 // Save word to memory
	MSTORE8 Opcode = 0x53 // Save byte to memory
	MSIZE   Opcode = 0x59 // Get size of active memory in bytes
	MCOPY   Opcode = 0x5E // Copy memory areas
)

// Storage opcodes
const (
	SLOAD  Opcode = 0x54 // Load word from storage
	SSTORE Opcode = 0x55 // Save word to storage
)

// TransientStorage opcodes
const (
	TLOAD  Opcode = 0x5C // Load word from transient storage
	TSTORE Opcode = 0x5D // Save word to transient storage
)

// Control flow opcodes
const (
	LOG0     Opcode = 0xA0 // Append log record with no topics
	LOG1     Opcode = 0xA1 // Append log record with one topic
	LOG2     Opcode = 0xA2 // Append log record with two topics
	LOG3     Opcode = 0xA3 // Append log record with three topics
	LOG4     Opcode = 0xA4 // Append log record with four topics
	JUMP     Opcode = 0x56 // Alter the program counter
	JUMPI    Opcode = 0x57 // Conditionally alter the program counter
	JUMPDEST Opcode = 0x5B // Mark a valid destination for jumps
)

const (
	CREATE       Opcode = 0xF0 // Create a new account with associated code
	CALL         Opcode = 0xF1 // Message-call into an account
	CALLCODE     Opcode = 0xF2 // Message-call with another account's code only
	RETURN       Opcode = 0xF3 // Halt execution returning output data
	DELEGATECALL Opcode = 0xF4 // Message-call into this account with an alternative account's code
	CREATE2      Opcode = 0xF5 // Create a new account with associated code at a predictable address
	STATICCALL   Opcode = 0xFA // Static message-call into an account
	REVERT       Opcode = 0xFD // Halt execution and revert state changes but return data and remaining gas
	SELFDESTRUCT Opcode = 0xFF // Halt execution and register account for later deletion
)

#!/bin/sh

outname=./sample.parquet

sample_jsonl(){
	jq -c -n '[
		{timestamp:"2025-10-01T00:11:50.012345Z", severity:"INFO", status:200},
		{timestamp:"2025-10-01T00:11:50.012345Z", severity:"WARN", status:500}
	]' |
		jq -c '.[]'
}

conv_native(){
	sample_jsonl |
		./jsons2parquet |
		cat > "${outname}"
}

conv_wasi(){
	sample_jsonl |
		wazero run ./logs2parquet.wasm |
		cat > "${outname}"
}

conv(){
	test -f ./logs2parquet.wasm && conv_wasi
	test -f ./logs2parquet.wasm || conv_native
}

conv

which parquet-read | fgrep -q parquet-read || exec sh -c '
	echo parquet-read missing.
	echo you can install it using cargo install parquet --bin ...
	exit 1
'

parquet-read "${outname}"

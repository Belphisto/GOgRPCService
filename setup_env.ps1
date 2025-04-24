[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

$env:GOPATH = "F:\Documents\GitHub\GOgRPCService\go-env"
$env:GOCACHE = "F:\Documents\GitHub\GOgRPCService\go-build"
$env:Path += ";F:\Documents\GitHub\GOgRPCService\go-env\bin"
$env:Path += ";F:\Documents\install\protoc-30.2-win64\bin"


Write-Host "GOPATH installed: $env:GOPATH"
Write-Host "GOCACHE installed: $env:GOCACHE"

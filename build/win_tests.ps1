<#
    .SYNOPSIS
        This script tests builds the New Relic Infrastructure Agent
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",

    # Skip tests
    [switch]$skipTests=$false
)

echo "--- Checking dependencies"

echo "Checking Go..."
go version
if (-not $?)
{
    echo "Can't find Go"
    exit -1
}

if (-Not $skipTests) {
    echo "--- Running tests"

    go test .\src\...
    if (-not $?)
    {
        echo "Failed running tests"
        exit -1
    }
}
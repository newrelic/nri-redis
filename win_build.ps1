<#
    .SYNOPSIS
        This script verifies, tests, builds and packages a New Relic Infrastructure Integration
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",
    [string]$tag="v0.0.0"
)

$integration = $(Split-Path -Leaf $PSScriptRoot)
$integrationName = $integration.Replace("nri-", "")
$executable = "nri-$integrationName.exe"

$version=$tag.substring(1)

# verifying version number format
$v = $version.Split(".")

if ($v.Length -ne 3) {
    echo "-version must follow a numeric major.minor.patch semantic versioning schema (received: $version)"
    exit -1
}

$wrong = $v | ? { (-Not [System.Int32]::TryParse($_, [ref]0)) -or ( $_.Length -eq 0) -or ([int]$_ -lt 0)} | % { 1 }
if ($wrong.Length  -ne 0) {
    echo "-version major, minor and patch must be valid positive integers (received: $version)"
    exit -1
}


echo "Checking MSBuild.exe..."
$msBuild = (Get-ItemProperty hklm:\software\Microsoft\MSBuild\ToolsVersions\4.0).MSBuildToolsPath
if ($msBuild.Length -eq 0) {
    echo "Can't find MSBuild tool. .NET Framework 4.0.x must be installed"
    exit -1
}
echo $msBuild

$env:GOOS="windows"
$env:GOARCH=$arch


echo "--- Building Installer"

Push-Location -Path "pkg\windows\nri-$arch-installer"
$env:integration = $integration
. $msBuild/MSBuild.exe nri-installer.wixproj

if (-not $?)
{
    echo "Failed building installer"
    Pop-Location
    exit -1
}

echo "Making versioned installed copy"


cd bin\Release

cp "$integration-$arch.msi" "$integration-$arch.$version.msi"
cp "$integration-$arch.msi" "$integration.msi"


Pop-Location
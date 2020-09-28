<#
    .SYNOPSIS
        This script creates the win .MSI
#>
param (
    # Target architecture: amd64 (default) or 386
    [ValidateSet("amd64", "386")]
    [string]$arch="amd64",
    [string]$tag="v0.0.0",
    [string]$pfx_certificate_base64="none",
    [string]$pfx_passphrase="none"
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

echo "===> Import .pfx certificate from GH Secrets"
Import-PfxCertificate -FilePath mycert.pfx -Password (ConvertTo-SecureString -String $pfx_passphrase -AsPlainText -Force) -CertStoreLocation Cert:\CurrentUser\My

echo "===> Show certificate installed"
Get-ChildItem -Path cert:\CurrentUser\My\

echo "===> Configuring version $version for artifacts"
.\windows_set_version.ps1 -major $v[0] -minor $v[1] -patch $v[2]

echo "===> Checking MSBuild.exe..."
$msBuild = (Get-ItemProperty hklm:\software\Microsoft\MSBuild\ToolsVersions\4.0).MSBuildToolsPath
if ($msBuild.Length -eq 0) {
    echo "Can't find MSBuild tool. .NET Framework 4.0.x must be installed"
    exit -1
}
echo $msBuild

$env:GOOS="windows"
$env:GOARCH=$arch

echo "===> Building Installer"
Push-Location -Path "pkg\windows\nri-$arch-installer"

echo "===================>debug"
pwd
echo "===1"
ls ..\..\..\
echo "===2"
ls ..\..\..\dist\
echo "===3"
ls ..\..\..\dist\nri-redis-win_windows_$arch\



$env:integration = $integration
. $msBuild/MSBuild.exe nri-installer.wixproj

if (-not $?)
{
    echo "Failed building installer"
    Pop-Location
    exit -1
}

echo "===> Making versioned installed copy"
cd bin\Release
cp "$integration-$arch.msi" "$integration-$arch.$version.msi"
cp "$integration-$arch.msi" "$integration.msi"

Pop-Location
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

$buildYear = (Get-Date).Year
$integrationName = "redis"

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
Import-PfxCertificate -FilePath wincert.pfx -Password (ConvertTo-SecureString -String $pfx_passphrase -AsPlainText -Force) -CertStoreLocation Cert:\CurrentUser\My

echo "===> Show certificate installed"
Get-ChildItem -Path cert:\CurrentUser\My\

echo "===> Checking MSBuild.exe..."
$msBuild = (Get-ItemProperty hklm:\software\Microsoft\MSBuild\ToolsVersions\4.0).MSBuildToolsPath
if ($msBuild.Length -eq 0) {
    echo "Can't find MSBuild tool. .NET Framework 4.0.x must be installed"
    exit -1
}
echo $msBuild

echo "===> Building Installer"
Push-Location -Path "build\package\windows\nri-$arch-installer"

. $msBuild/MSBuild.exe nri-installer.wixproj /p:IntegrationVersion=${version} /p:IntegrationName=$integrationName /p:Year=$buildYear

if (-not $?)
{
    echo "Failed building installer"
    Pop-Location
    exit -1
}

echo "===> Making versioned installed copy"
cd bin\Release
cp "nri-$integrationName-$arch.msi" "nri-$integrationName-$arch.$version.msi"
# todo: why do we need this?
cp "nri-$integrationName-$arch.msi" "nri-$integrationName.msi"

Pop-Location
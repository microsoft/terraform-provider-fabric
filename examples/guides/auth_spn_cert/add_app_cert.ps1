$certValue = [System.Convert]::ToBase64String([System.IO.File]::ReadAllBytes('client.crt'))
Get-AzADApplication -ApplicationId '00000000-0000-0000-0000-000000000000' | New-AzADAppCredential -CertValue $certValue

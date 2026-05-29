---
version: v1.13.0
---

# Using hiera-eyaml as backend for secrets' encryption

Helmsman uses helm-secrets as a default solution for secrets' encryption.
And while it is a good off-the-shelve solution it may quickly start causing problems when few developers start working on the secrets files simultaneously.
SOPS-based secrets can not be easily merged or rebased in case of conflicts etc.
That is why another solution for secrets organised in YAMLs was proposed in [hiera-eyaml](https://github.com/voxpupuli/hiera-eyaml).

## Example

Having environment defined with:

* example.yaml:
```yaml
settings:
  eyamlEnabled: true
```

Helmsman will use hiera-eyaml gem to decrypt secrets files defined for applications.
They public and private keys should be placed in `keys` directory with names of `public_key.pkcs7.pem` and `private_key.pkcs7.pem`.
The keys' path can be overwritten with

```yaml
settings:
  eyamlEnabled: true
  eyamlPrivateKeyPath: ../keys/custom.pem
  eyamlPublicKeyPath: ../keys/custom.pub
```

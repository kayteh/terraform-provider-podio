# Oidop Terraform Provider

This is a not the official Podio provider. This provider follows a very small set of goals, please do not use it in production.

It is tied to an unofficial podio-go SDK as well, maintained by only myself.

## API Key "Trust" Levels

A major caveat to this whole system is the use of "trust" levels to moderate API access. This terraform provider attempts to make certain operations that require higher trust break out gracefully. The default behavior is to not do this.

You need a trust level of 2 (defaulting to 0) to use this to it's full potential. Contact Podio support and let them know you are using this tool and you can link them to this documentation as to why.

This Terraform provider needs to:

- Create, modify, and delete spaces
- Invite and remove members from spaces

## Current Target: MMF1: Kanban

- [x] Workspace creation
- [x] App creation
- [ ] Template/Field creation (Text + Category)
- [ ] Icon search picker
- [ ] Item creation

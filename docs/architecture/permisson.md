# Echelon Permissions

The `pkg/core/permissions` package is responsible for calculating and
checking permissions.

The package contains only permission-domain logic and does not directly
depend on MongoDB, Redis, JWT, HTTP, or Gin.

---

## Module Structure

```text
pkg/core/permissions/
├── permission.go
├── mask.go
├── permission_value.go
├── override.go
├── default.go
├── input.go
├── roles.go
├── server.go
├── channel.go
├── hierarchy.go
├── delegation.go
│
├── server_test.go
├── channel_test.go
├── hierarchy_test.go
└── delegation_test.go
```

---

# 1. `permission.go`

Defines the individual permission capabilities.

## `type Permission uint64`

Each permission is represented by one bit.

Example:

```go
SendMessage Permission = 1 << 22
```

This means `SendMessage` owns bit `22`.

## Permission Constants

| Permission | Bit |
|---|---:|
| ManageChannel | 0 |
| ManageServer | 1 |
| ManagePermissions | 2 |
| ManageRole | 3 |
| ManageCustomisation | 4 |
| KickMembers | 6 |
| BanMembers | 7 |
| TimeoutMembers | 8 |
| AssignRoles | 9 |
| ChangeNickname | 10 |
| ManageNicknames | 11 |
| ChangeAvatar | 12 |
| RemoveAvatars | 13 |
| ViewChannel | 20 |
| ReadMessageHistory | 21 |
| SendMessage | 22 |
| ManageMessages | 23 |
| ManageWebhooks | 24 |
| InviteOthers | 25 |
| SendEmbeds | 26 |
| UploadFiles | 27 |
| Masquerade | 28 |
| React | 29 |
| Connect | 30 |
| Speak | 31 |
| Video | 32 |
| MuteMembers | 33 |
| DeafenMembers | 34 |
| MoveMembers | 35 |
| Listen | 36 |
| MentionEveryone | 37 |
| MentionRoles | 38 |
| BypassSlowmode | 39 |
| ViewAuditLogs | 40 |

Some bit positions are reserved for future permissions.

Normal permissions use bits `0–51`.

Permission bit positions must never be renumbered or reused.

## `GrantAllSafe`

Represents all normal assignable permission bits.

It is used when a privileged user or server owner receives all normal
permissions.

---

# 2. `mask.go`

Defines how multiple permissions are stored together.

## `type Mask uint64`

A `Mask` represents a set of permissions.

Example:

```go
mask := Mask(ViewChannel) |
       Mask(ReadMessageHistory) |
       Mask(SendMessage)
```

The resulting mask contains all three permissions.

## `Has(permission Permission) bool`

Checks whether a permission exists in the mask.

```go
mask.Has(SendMessage)
```

Returns:

```text
true  → permission exists
false → permission does not exist
```

## `HasAll(required Mask) bool`

Checks whether all permissions in another mask exist.

```go
mask.HasAll(required)
```

This is used when checking whether an actor has every permission they
are trying to delegate.

## Internal Mutation Operations

The mask also contains internal operations for modifying the bitset:

```text
add()
remove()
clear()
grantAllSafe()
```

These are implementation details and are not required by callers.

---

# 3. `permission_value.go`

Represents the final calculated permissions.

## `type PermissionValue`

`PermissionValue` contains the effective permissions calculated for a
user.

The internal mask is not exposed directly.

Callers check permissions through methods such as:

```go
value.Has(SendMessage)
```

## `FromRaw(raw uint64)`

Creates a `PermissionValue` from a raw permission value.

This is primarily used internally when constructing permission results.

## `IntoRaw() uint64`

Converts the permission value back into its raw `uint64`
representation.

This can be used when another layer needs to store or transmit the
permission value.

## `Allow(permission Permission)`

Internally adds a permission to the current value.

## `Revoke(permission Permission)`

Internally removes a permission from the current value.

## `RevokeAll()`

Internally removes all permissions.

## `Restrict(mask Mask)`

Restricts the current permissions to the permissions contained in the
provided mask.

Example:

```text
Current:
ViewChannel
ReadMessageHistory
SendMessage
ManageMessages

Restriction:
ViewChannel
ReadMessageHistory

Result:
ViewChannel
ReadMessageHistory
```

This is used for timeout restrictions.

## `Has(permission Permission) bool`

Checks whether the final permission value contains a permission.

This is one of the main public APIs of the package.

## `HasAll(required Mask) bool`

Checks whether the final permission value contains all permissions in
the provided mask.

---

# 4. `override.go`

Defines permission overrides.

## `type Override`

```go
type Override struct {
    Allow Mask
    Deny  Mask
}
```

An override can:

```text
Allow permissions
Deny permissions
```

Example:

```text
Allow:
    SendMessage

Deny:
    ManageMessages
```

## `Validate()`

Checks that the same permission does not exist in both `Allow` and
`Deny`.

Invalid:

```text
Allow:
    SendMessage

Deny:
    SendMessage
```

## `Apply(current PermissionValue)`

Applies the override to an existing permission value.

The operation is:

```text
1. Add allowed permissions
2. Remove denied permissions
```

This operation is part of the internal permission calculation flow.

---

# 5. `default.go`

Contains predefined permission masks.

## `DEFAULT_PERMISSION_VIEW_ONLY`

Contains:

```text
ViewChannel
ReadMessageHistory
```

Used when a user should only be able to view a channel.

## `DEFAULT_PERMISSION`

Normal default permissions:

```text
ViewChannel
ReadMessageHistory
SendMessage
InviteOthers
SendEmbeds
UploadFiles
Connect
Speak
Listen
Video
```

## `DEFAULT_PERMISSION_SERVER`

Server default permissions extend the normal defaults with:

```text
React
ChangeNickname
ChangeAvatar
```

## `ALLOW_IN_TIMEOUT`

Permissions that remain available while a user is timed out:

```text
ViewChannel
ReadMessageHistory
```

---

# 6. `input.go`

Contains the input structures used by the permission engine.

The permission package does not query a database itself.

Instead, the caller gathers the required authorization information and
passes it through these structures.

## `RolePermissionInput`

Represents a role and its permission override.

```go
type RolePermissionInput struct {
    ID       string
    Rank     int
    Override Override
}
```

Fields:

```text
ID       → identifies the role
Rank     → role authority level
Override → permissions allowed/denied by the role
```

## `ServerPermissionInput`

Contains the information required to calculate server permissions.

Important fields:

```text
IsPrivileged
IsOwner
IsMember
DefaultPermissions
Roles
TimedOut
CanPublish
CanReceive
```

## `ChannelPermissionInput`

Contains the information required to calculate channel permissions.

Important fields:

```text
ServerPermissions
DefaultOverride
RoleOverrides
MemberOverride
CanPublish
CanReceive
TimedOut
```

`ServerPermissions` is already calculated by `ResolveServer()`.

---

# 7. `roles.go`

Contains role ordering logic.

## `SortRoles()`

Sorts roles into deterministic authority order.

The ordering is:

```text
Higher numeric rank → evaluated first
Lower numeric rank   → evaluated later
```

Example:

```text
Rank 10
Rank 5
Rank 2
```

is evaluated as:

```text
10 → 5 → 2
```

The higher-authority role at rank `2` is therefore applied later and can
override conflicting permissions.

The function also:

```text
rejects duplicate role IDs
does not modify the original input slice
uses role ID as a deterministic tie-breaker
```

The role ordering is important because multiple roles can contain
conflicting permission overrides.

---

# 8. `server.go`

Calculates effective permissions for a user at the server level.

## `ResolveServer(input ServerPermissionInput) PermissionValue`

This is one of the main public APIs.

The evaluation order is:

```text
1. Privileged user
2. Server owner
3. Non-member
4. Default server permissions
5. Role overrides
6. Publish restriction
7. Receive restriction
8. Timeout restriction
9. Final permissions
```

### Privileged User

A privileged user receives:

```text
GrantAllSafe
```

### Server Owner

The server owner receives:

```text
GrantAllSafe
```

### Non-member

A user who is not a server member receives:

```text
0
```

### Normal Member

A normal member starts with:

```text
DefaultPermissions
```

Role overrides are then applied according to role rank.

### Publish Restriction

If:

```text
CanPublish = false
```

the following permissions are removed:

```text
Speak
Video
```

### Receive Restriction

If:

```text
CanReceive = false
```

the following permission is removed:

```text
Listen
```

### Timeout Restriction

If:

```text
TimedOut = true
```

permissions are restricted to:

```text
ViewChannel
ReadMessageHistory
```

---

# 9. `channel.go`

Calculates effective permissions inside a specific channel.

## `ResolveChannel(input ChannelPermissionInput) PermissionValue`

This is the second main public API.

Channel calculation starts with:

```text
ServerPermissions
```

and then applies channel-specific rules.

The evaluation order is:

```text
1. Server permissions
2. Channel default / @everyone override
3. Channel role overrides
4. Member-specific override
5. Publish restriction
6. Receive restriction
7. Timeout restriction
8. ViewChannel gate
9. Final permissions
```

## Channel Default Override

Applies the channel's `@everyone` override.

## Channel Role Overrides

Applies overrides belonging to the user's roles.

Roles are evaluated according to the same deterministic role ordering
used by server permission calculation.

## Member Override

Applies an override specifically for the user.

This allows a channel to grant or deny permissions for one specific
member.

## Publish Restriction

If publishing is disabled:

```text
Speak
Video
```

are removed.

## Receive Restriction

If receiving is disabled:

```text
Listen
```

is removed.

## Timeout Restriction

If the user is timed out, the result is restricted to:

```text
ViewChannel
ReadMessageHistory
```

## ViewChannel Gate

After all channel permissions are calculated, the engine checks:

```text
ViewChannel
```

If the permission is not present, the function returns:

```text
0
```

This prevents a user from interacting with a channel that they cannot
view.

---

# 10. `hierarchy.go`

Contains role hierarchy checks.

## `CanManageRole(actorRank, targetRank int) bool`

Determines whether an actor can manage a target role.

The rule is:

```text
actorRank < targetRank
```

Example:

```text
Actor rank: 5
Target rank: 10

5 < 10

Allowed
```

Equal rank:

```text
Actor rank: 5
Target rank: 5

5 < 5

Denied
```

Higher-authority target:

```text
Actor rank: 5
Target rank: 2

5 < 2

Denied
```

The owner is handled separately by the caller because ownership is not
represented by the function arguments.

## `CanManageMember(actorRank, targetRank int) bool`

Uses the same hierarchy rule:

```text
actorRank < targetRank
```

The actor must have strictly higher authority than the target member.

---

# 11. `delegation.go`

Controls permission delegation.

## `CanDelegatePermissions(input PermissionDelegationInput) bool`

Determines whether an actor can grant a requested set of permissions to
a target role.

The function checks:

```text
1. Is the actor the owner?
2. Can the actor manage the target role?
3. Does the actor have all permissions they are trying to grant?
```

An actor cannot grant permissions outside their own authority.

Example:

```text
Actor permissions:
    SendMessage

Requested permissions:
    SendMessage
```

Allowed if the target role is manageable.

But:

```text
Actor permissions:
    SendMessage

Requested permissions:
    ManageServer
```

Denied.

The actor also cannot bypass role hierarchy.

The owner can manage normally.

---

# 12. Permission Calculation Flow

## Server Permission Calculation

```text
ServerPermissionInput
        ↓
ResolveServer()
        ↓
PermissionValue
```

## Channel Permission Calculation

```text
ServerPermissionInput
        ↓
ResolveServer()
        ↓
Server PermissionValue
        ↓
ChannelPermissionInput
        ↓
ResolveChannel()
        ↓
Channel PermissionValue
```

The resulting permission value can then be checked:

```go
if value.Has(SendMessage) {
    // allowed
}
```

---

# 13. Complete Example

Consider a server with the following permissions.

## @everyone

```text
ViewChannel
ReadMessageHistory
SendMessage
```

## Moderator

The Moderator role adds:

```text
ManageMessages
```

Therefore a Moderator receives:

```text
ViewChannel
ReadMessageHistory
SendMessage
ManageMessages
```

## Alice

Alice has the:

```text
Moderator
```

role.

Therefore Alice's server permissions are:

```text
ViewChannel
ReadMessageHistory
SendMessage
ManageMessages
```

## #staff

The `#staff` channel has an `@everyone` override:

```text
Deny:
    SendMessage
```

Channel evaluation:

```text
Server permissions
        ↓
ViewChannel
ReadMessageHistory
SendMessage
ManageMessages
        ↓
#staff @everyone override
        ↓
Deny SendMessage
        ↓
Final permissions
```

Alice's final permissions in `#staff` are:

```text
ViewChannel
ReadMessageHistory
ManageMessages
```

Therefore:

```text
Alice can view #staff
Alice can read message history
Alice can manage messages
Alice cannot send messages
```

The channel override changes Alice's effective permissions for
`#staff` without changing her server-level permissions.

---

# 14. Test Files

The permission module has tests for each major component.

## `server_test.go`

Tests server permission calculation:

```text
Owner
Privileged user
Non-member
Normal member
Role overrides
Multiple roles
Conflicting roles
Timeout
Publish restriction
Receive restriction
Role ordering
```

## `channel_test.go`

Tests channel permission calculation:

```text
Server permission inheritance
@everyone override
Channel role override
Member override
Override precedence
Timeout
Publish restriction
Receive restriction
ViewChannel gate
```

## `hierarchy_test.go`

Tests role hierarchy:

```text
Actor above target
Actor equal to target
Actor below target
Multiple rank combinations
```

## `delegation_test.go`

Tests permission delegation:

```text
Valid delegation
Invalid delegation
Grant beyond authority
Grant to higher role
Grant to equal role
Owner delegation
```

---

# 15. Package Boundary

The permission package is intentionally independent from infrastructure.

It does not directly access:

```text
MongoDB
Redis
JWT
Gin
HTTP
```

Instead, the surrounding application gathers the required information
and creates permission inputs.

The flow is:

```text
HTTP / Gin
     ↓
Authentication
     ↓
Authorization Service
     ↓
Database / other data sources
     ↓
Permission Input
     ↓
Permission Engine
     ↓
PermissionValue
     ↓
Has(permission)
```

The permission engine itself only works with permission-domain types and
values.

This keeps permission calculation deterministic, reusable, and
independent of infrastructure.

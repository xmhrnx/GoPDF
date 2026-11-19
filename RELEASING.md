# Releasing GoPDF

Publishing a GitHub release with a tag like `v0.1.0` runs
`.github/workflows/release.yml`. The tag sets the version and the run number
the build number. The workflow tests, builds every package and attaches them
to the release. It is written by the Giode CLI from `giode.toml`: after
changing `[ci]` or a store section, run `giode ci -f` instead of editing it,
then redo the two edits the comment in `[ci]` lists (the pinned giode CLI
version and `go-version-file: go.mod`).

| Where | What the release gets | Needs |
| --- | --- | --- |
| GitHub, Linux | `.deb` and `.rpm`, x64 | nothing |
| GitHub, Windows | `setup.exe` installer, x64 | nothing; a code signing certificate stops SmartScreen warnings |
| GitHub, macOS | a universal `.dmg` | a Developer ID certificate, or macOS refuses to open it |
| Mac App Store | a `.pkg` to upload with Transporter | an Apple Developer account and the store certificates |
| Microsoft Store | an `.msix` package to upload in Partner Center | a Partner Center account and the reserved name |

There is no Google Play or iOS App Store package: GoPDF is a desktop app (a
frameless, keyboard-driven window with menus, a 900 × 560 minimum size), and
MuPDF, which renders the pages, has no iOS build in go-fitz. A phone version
would be a different app, not a store package of this one.

The Windows packages are built on a Windows runner in MSYS2's MINGW64
(`[ci] windows_msys2`): MuPDF is linked through cgo, and cross-building on
Linux would turn cgo off and leave the program needing a MuPDF DLL it doesn't
ship.

## Repository secrets

Set these in GitHub under Settings → Secrets and variables → Actions. Each
`*_BASE64` secret is the file encoded with `base64 -i file | pbcopy` on macOS,
or `base64 -w0 file` on Linux. Without a secret, that step is skipped and the
rest of the release still builds.

| Secret | What it is |
| --- | --- |
| `MACOS_CERT_BASE64`, `MACOS_CERT_PASS` | The "Developer ID Application" certificate as a `.p12`, and its password: signs the `.dmg` |
| `GIODE_NOTARY_ID`, `GIODE_NOTARY_TEAM_ID`, `GIODE_NOTARY_PASS` | Your Apple ID, team ID and an app-specific password: notarizes the `.dmg` |
| `MACOS_STORE_CERT_BASE64`, `MACOS_STORE_CERT_PASS` | The "Apple Distribution" certificate as a `.p12` |
| `MACOS_INSTALLER_CERT_BASE64`, `MACOS_INSTALLER_CERT_PASS` | The "Mac Installer Distribution" certificate as a `.p12` |
| `MACOS_PROFILE_BASE64` | The Mac App Store provisioning profile (`.provisionprofile`) |
| `WINDOWS_CERT_BASE64`, `GIODE_WINDOWS_CERT_PASS` | An Authenticode `.pfx` and its password: signs the Windows installer (optional) |

The Mac App Store `.pkg` is built only when all five `MACOS_STORE_*`,
`MACOS_INSTALLER_*` and `MACOS_PROFILE_BASE64` secrets are set. The MSIX needs
no secret: the Microsoft Store signs it.

## Mac App Store

1. Join the Apple Developer Program.
2. In Certificates, Identifiers & Profiles:
   - Register the App ID `com.xmhrnx.gopdf` (the `app.id` in `giode.toml`).
     It needs no capabilities.
   - Create an "Apple Distribution" and a "Mac Installer Distribution"
     certificate, and export each from Keychain Access as a `.p12`.
   - Create a "Mac App Store Connect" provisioning profile for the App ID.
   - For the `.dmg`: a "Developer ID Application" certificate, and an
     app-specific password for your Apple ID (account.apple.com).
3. In App Store Connect, create a macOS app with the bundle ID
   `com.xmhrnx.gopdf`. If the name "GoPDF" is taken there, use a longer one
   such as "GoPDF – PDF Editor" (the listing name only; the app keeps its own).
4. Add the secrets above, then publish a release. The release gets
   `GoPDF-<version>.pkg`; upload it with the Transporter app.
5. In App Store Connect, fill in the listing (see "Store listing" below):
   - **App Privacy:** "Data Not Collected". GoPDF makes no network
     connections. The build's privacy manifest says the same
     (`[privacy] tracking = false`, nothing collected).
   - **Privacy policy URL:** https://xmhrnx.github.io/GoPDF/#privacy
   - **Category:** Productivity (already set in the build).
   - **Age rating:** 4+: no objectionable content of its own, no web access,
     no user-to-user communication.
   - **Export compliance:** answered in the build
     (`ITSAppUsesNonExemptEncryption = false`). The only encryption is
     standard algorithms from Go's standard library: AES-256 when the user
     password-protects a PDF, and RSA/SHA-256 for digital signatures.

**The App Sandbox limits GoPDF in the store version.** The build has only the
`files.user-selected.read-write` entitlement (no network, which GoPDF never
uses):

- **Files:** GoPDF can open the files the user picks in the open panel, drops
  on the window or opens from Finder, save them in place, and write wherever
  the save panel or Split's folder panel points. Saving in place writes the
  file directly, since the sandbox doesn't let it put a temporary file next to
  it.
- **Recent files:** the sandbox grants a file for the session it was opened
  in. After a relaunch, entries on the welcome screen's recent list may not
  open (permission denied) until the file is opened again with File › Open or
  from Finder.
- **Settings and digital IDs** are kept in the app's container
  (`~/Library/Containers/com.xmhrnx.gopdf/Data/Library/Application Support`),
  separate from the `.dmg` version's.

The notarized `.dmg` on GitHub has none of these limits.

## Microsoft Store

1. In Partner Center (the same account as EasySend and AnyFile), reserve the
   name "GoPDF". If it's taken, reserve another and set `app.name` for the
   MSIX build to it, as AnyFile's release does.
2. `giode.toml` already has the account's values under `[windows.msix]`:
   `publisher = "CN=B40E8F57-A7E0-40BB-8E74-5D051D0B578B"` and
   `publisher_display_name = "Mehran Zahiroddini"`. `identity_name` is set to
   `MehranZahir.GoPDF`: once the name is reserved, open the product's Product
   identity page and check that **Package/Identity/Name** matches it exactly;
   if not, change it and commit before the release.
3. Publish a release. It gets `GoPDF-<version>-amd64.msix`; upload it in a
   Partner Center submission. The Store signs it, so it needs no certificate
   here. The package registers GoPDF for `.pdf` files ("Open with").
4. GoPDF is a full-trust desktop app (a Win32 program that reads and writes
   the PDFs the user opens). When the submission asks why it uses the
   `runFullTrust` capability, say that.
5. **Privacy policy URL:** https://xmhrnx.github.io/GoPDF/#privacy. In the
   questionnaire, GoPDF doesn't access, collect or transmit personal
   information. **Age rating (IARC):** no violence, sexual content, gambling,
   user interaction, location sharing or purchases: Everyone / 3+.

## Store listing

**Name:** GoPDF

**Subtitle (Apple, ≤ 30 characters):** Annotate, sign and merge PDFs

**Short description (≤ 80 characters):**
Annotate, sign, merge, split and protect PDFs. Native and fast, fully offline.

**Keywords (Apple, ≤ 100 characters):**
pdf,editor,annotate,sign,signature,digital signature,merge,split,highlight,watermark,encrypt,pages

**Description:**

> GoPDF is a fast, native PDF editor. There is no browser engine inside: it is
> written in Go and drawn straight to the GPU, and only the pages in view are
> rendered, so even long documents scroll smoothly.
>
> Annotate: text boxes, pen, marker, highlight, rectangles, ellipses, lines,
> arrows and whiteout, with your own colors, widths and font sizes. Select
> markup to move, restyle, copy or delete it. When you save, the markup is
> drawn into the page, so every PDF viewer shows it.
>
> Sign: draw your signature or type your name in one of six script
> typefaces, keep it, and place it in any document. Turn it into a digital
> signature with a digital ID you create in GoPDF or import from a
> certificate authority (.p12/.pfx). GoPDF checks the signatures in the files
> you open and shows who signed and whether the document changed since.
>
> Organize: rotate, duplicate, reorder, insert and delete pages. Merge PDFs,
> extract page ranges, and split a document into parts.
>
> Stamp text watermarks and page numbers, password-protect files with
> AES-256, edit the title and author, and export a page as PNG or the text as
> a .txt file. Undo and redo, light and dark themes, and keyboard shortcuts
> for everything.
>
> Private by design: GoPDF never connects to the internet. Your documents,
> signatures and digital IDs stay on your computer.

**Category:** Productivity (Mac App Store), Productivity (Microsoft Store).

**Screenshots:** the website draws its page illustration in HTML, so there
are no screenshot files yet. Capture the app in light and dark mode with a
document open: annotating (markup selected, the inspector showing its
style), the signature dialog, a signed document with the signature details,
and the page thumbnails while reordering. Mac 2880 × 1800 (or 1440 × 900),
Microsoft Store 1920 × 1080.

## Licensing before submitting

GoPDF is Apache-2.0, but every package links MuPDF (through go-fitz), which
is AGPL-3.0, so the program as distributed falls under the AGPL. The source
is public, which the AGPL asks for. The Mac App Store's terms of use are
widely considered to conflict with the GPL family (VLC was pulled for this),
and Artifex, MuPDF's owner, sells commercial licenses for app store
distribution. Decide this, or ask Artifex, before the Mac App Store
submission. The Microsoft Store's terms don't have this problem.

## Linux

The `.deb` and `.rpm` are attached to the release, and the website's download
button links to the release. They list the system libraries they need, so
`apt` and `dnf` install those too.

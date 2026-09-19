---@meta x509

---@class x509.Certificate
---@field dns_names string[] SAN dNSName entries
---@field email_addresses string[] SAN rfc822Name entries
---@field ext_key_usage string[] e.g. {"server_auth", "client_auth"}
---@field fingerprint_sha256 string lowercase hex SHA-256 of the raw DER
---@field ip_addresses string[] SAN iPAddress entries, v4 and v6, in net.IP.String() form
---@field is_ca boolean basic constraints CA flag
---@field issuer string RFC 2253-style distinguished name
---@field issuer_cn string issuer's CommonName, "" if absent
---@field issuer_org string[] issuer's Organization values
---@field key_usage string[] e.g. {"digital_signature", "key_encipherment"}
---@field max_path_len number basic constraints path length constraint, 0 if unset
---@field not_after number Unix seconds
---@field not_before number Unix seconds
---@field public_key_algorithm string e.g. "RSA", "ECDSA", "Ed25519"
---@field public_key_bits number e.g. 2048 for a 2048-bit RSA key
---@field serial string exact decimal string; serials can exceed 2^53
---@field signature_algorithm string e.g. "SHA256-RSA"
---@field subject string RFC 2253-style distinguished name
---@field subject_cn string subject's CommonName, "" if absent
---@field subject_org string[] subject's Organization values
---@field uris string[] SAN uniformResourceIdentifier entries
---@field version number X.509 version (3 for a v3 certificate)

---@class x509.VerifyChainOptions
---@field at_time number Unix seconds; REQUIRED, must not be 0 — pass time.now()
---@field dns_name string hostname to check the leaf against; "" skips the hostname check
---@field key_usages string[] required extended key usages; empty means {"server_auth"}

---@class x509
local x509 = {}

--- parses the first CERTIFICATE PEM block in a string and returns its fields; raises on missing/wrong-type PEM block or bad DER
---@param pem string PEM-encoded certificate text
---@return x509.Certificate cert the parsed certificate
function x509.parse(pem) end

--- parses every CERTIFICATE PEM block in a string, in file order; raises if any CERTIFICATE block fails to parse or none are found
---@param pem string PEM text containing one or more CERTIFICATE blocks (other block types are ignored)
---@return x509.Certificate[] certs the parsed certificates, in the order they appear
function x509.parse_chain(pem) end

--- returns the number of days between now and the certificate's not_after; may be negative (already expired) or fractional
---@param pem string PEM-encoded certificate text
---@param now number the current time as Unix seconds, e.g. time.now()
---@return number days fractional days until expiry; negative if already expired
function x509.expires_in_days(pem, now) end

--- reports whether when falls within [not_before, not_after], inclusive of both bounds
---@param pem string PEM-encoded certificate text
---@param when number the time to check, as Unix seconds
---@return boolean ok true if not_before <= when <= not_after
function x509.is_valid_at(pem, when) end

--- verifies leaf against intermediates and roots; NEVER consults the system trust store, so an empty roots argument can never succeed
---@param leaf string PEM-encoded leaf certificate to verify
---@param intermediates string PEM text with zero or more intermediate CA certificates
---@param roots string PEM text with zero or more trusted root certificates; empty means nothing is trusted
---@param opts x509.VerifyChainOptions required options: dns_name, at_time (Unix seconds, required, must not be 0), key_usages (empty defaults to {"server_auth"})
---@return boolean ok true if the chain verifies against roots
---@return string reason empty string on success; the underlying x509 verification error text on failure
function x509.verify_chain(leaf, intermediates, roots, opts) end

return x509

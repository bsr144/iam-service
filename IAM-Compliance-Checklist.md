# IAM System - Compliance Checklist
# ISO 27001 and PCI DSS Alignment

---

## Document Control

| Attribute | Value |
|-----------|-------|
| **Document Title** | IAM System Compliance Checklist |
| **Version** | 1.0 |
| **Created Date** | 2026-02-03 |
| **Last Updated** | 2026-02-03 |
| **Parent Document** | IAM-Implementation-Plan-V2.md |
| **Review Frequency** | Quarterly |

---

## 1. Executive Summary

This document provides a comprehensive compliance checklist for the IAM System, covering:
- **ISO 27001:2022** - Information Security Management System
- **PCI DSS 4.0** - Payment Card Industry Data Security Standard (where applicable)
- **OWASP Top 10 2021** - Web Application Security Risks
- **GDPR** - General Data Protection Regulation (for personal data handling)

### 1.1 Compliance Status Legend

| Status | Symbol | Description |
|--------|--------|-------------|
| Compliant | ✅ | Fully implemented and verified |
| Partial | 🔶 | Partially implemented, needs work |
| Non-Compliant | ❌ | Not implemented |
| Not Applicable | N/A | Requirement not applicable to IAM |
| Planned | 📋 | Scheduled for implementation |

---

## 2. ISO 27001:2022 Compliance Checklist

### 2.1 Access Control (A.9)

#### A.9.1 Business Requirements of Access Control

| Control | Requirement | IAM Implementation | Status | Evidence |
|---------|-------------|-------------------|--------|----------|
| A.9.1.1 | Access control policy | RBAC policy defined; permission model documented | 📋 | PRD-Section2-F3-Authorization.md |
| A.9.1.2 | Access to networks and network services | API authentication required; JWT validation | 📋 | TRD-Section5-Security.md |

#### A.9.2 User Access Management

| Control | Requirement | IAM Implementation | Status | Evidence |
|---------|-------------|-------------------|--------|----------|
| A.9.2.1 | User registration and de-registration | Self-registration with approval workflow; admin user creation; status lifecycle (ACTIVE→INACTIVE) | 📋 | PRD-Section2-F2-UserManagement.md |
| A.9.2.2 | User access provisioning | Role assignment API; branch-scoped permissions | 📋 | PRD-Section2-F3-Authorization.md |
| A.9.2.3 | Management of privileged access rights | Platform admin role; tenant admin role; privilege separation | 📋 | TRD-Section3-DatabaseDesign.md |
| A.9.2.4 | Management of secret authentication information | Password policy; bcrypt hashing (cost 12); Vault for secrets | 📋 | TRD-Section5-Security.md |
| A.9.2.5 | Review of user access rights | Audit log for role assignments; GET /users/{id}/roles API | 📋 | PRD-Section2-F5-AuditLogging.md |
| A.9.2.6 | Removal or adjustment of access rights | Status lifecycle; role revocation API; soft delete | 📋 | PRD-Section2-F2-UserManagement.md |

#### A.9.3 User Responsibilities

| Control | Requirement | IAM Implementation | Status | Evidence |
|---------|-------------|-------------------|--------|----------|
| A.9.3.1 | Use of secret authentication information | Password complexity policy; PIN policy; no password sharing | 📋 | PRD Section 5 (Validation Rules) |

#### A.9.4 System and Application Access Control

| Control | Requirement | IAM Implementation | Status | Evidence |
|---------|-------------|-------------------|--------|----------|
| A.9.4.1 | Information access restriction | Permission check API; branch scoping; tenant isolation | 📋 | PRD-Section2-F3-Authorization.md |
| A.9.4.2 | Secure log-on procedures | Multi-step login (password + PIN); rate limiting; lockout | 📋 | PRD-Section2-F1-Authentication.md |
| A.9.4.3 | Password management system | Password history; complexity rules; change password API | 📋 | TRD-Section3-DatabaseDesign.md |
| A.9.4.4 | Use of privileged utility programs | Not applicable (API-only system) | N/A | - |
| A.9.4.5 | Access control to program source code | Not in scope (handled by Git/CI) | N/A | - |

### 2.2 Cryptography (A.10)

| Control | Requirement | IAM Implementation | Status | Evidence |
|---------|-------------|-------------------|--------|----------|
| A.10.1.1 | Policy on use of cryptographic controls | Encryption strategy documented | 📋 | TRD-Section5-Security.md |
| A.10.1.2 | Key management | JWT key rotation (30 days); Vault for key storage | 📋 | TRD-Section5-Security.md |

### 2.3 Operations Security (A.12)

#### A.12.4 Logging and Monitoring

| Control | Requirement | IAM Implementation | Status | Evidence |
|---------|-------------|-------------------|--------|----------|
| A.12.4.1 | Event logging | All auth/authz events logged to OpenSearch | 📋 | PRD-Section2-F5-AuditLogging.md |
| A.12.4.2 | Protection of log information | Immutable logs; access control; encrypted storage | 📋 | TRD-Section5-Security.md |
| A.12.4.3 | Administrator and operator logs | Admin actions logged with before/after state | 📋 | PRD-Section2-F5-AuditLogging.md |
| A.12.4.4 | Clock synchronization | All timestamps in UTC; NTP synchronized | 📋 | TRD-Section3-DatabaseDesign.md |

### 2.4 Compliance (A.18)

| Control | Requirement | IAM Implementation | Status | Evidence |
|---------|-------------|-------------------|--------|----------|
| A.18.1.3 | Protection of records | Audit log retention (90 days hot, 7 years archive) | 📋 | PRD-Section2-F5-AuditLogging.md |
| A.18.1.4 | Privacy and protection of PII | Profile data protection; consent tracking | 📋 | TRD-Section5-Security.md |

---

## 3. PCI DSS 4.0 Compliance Checklist

> **Note:** PCI DSS applies if the IAM system handles, processes, or stores cardholder data. If IAM only manages authentication for systems that handle card data, these controls ensure the authentication system itself is secure.

### 3.1 Requirement 7: Restrict Access to System Components

| Requirement | Control | IAM Implementation | Status | Evidence |
|-------------|---------|-------------------|--------|----------|
| 7.1.1 | Access control policy | RBAC policy; principle of least privilege | 📋 | PRD-Section2-F3-Authorization.md |
| 7.2.1 | Access based on job function | Role-based access; application-scoped permissions | 📋 | TRD-Section3-DatabaseDesign.md |
| 7.2.2 | Approval required for access | Approval workflow for user registration | 📋 | PRD-Section2-F2-UserManagement.md |
| 7.2.3 | Default deny | Permissions required; no implicit access | 📋 | PRD-Section2-F3-Authorization.md |
| 7.2.4 | All access changes documented | Audit log for role assignments | 📋 | PRD-Section2-F5-AuditLogging.md |
| 7.2.5 | Access provisioned by authorized personnel | Tenant admin creates users; platform admin creates tenants | 📋 | PRD-Section2-F4-OrganizationManagement.md |
| 7.2.6 | Access reviewed periodically | Audit query API for access review | 📋 | PRD-Section2-F5-AuditLogging.md |

### 3.2 Requirement 8: Identify Users and Authenticate Access

| Requirement | Control | IAM Implementation | Status | Evidence |
|-------------|---------|-------------------|--------|----------|
| 8.1.1 | Unique user IDs | UUID per user; email unique within tenant | 📋 | TRD-Section3-DatabaseDesign.md |
| 8.1.2 | Manage user identities | User lifecycle management; status tracking | 📋 | PRD-Section2-F2-UserManagement.md |
| 8.2.1 | Proper authentication | Password + PIN multi-factor | 📋 | PRD-Section2-F1-Authentication.md |
| 8.2.2 | Unique authentication | JTI in JWT; no shared credentials | 📋 | TRD-Section5-Security.md |
| 8.2.3 | Invalid authentication locked | Account lockout after 5 failures | 📋 | PRD-Section2-F1-Authentication.md |
| 8.2.4 | Session timeout | JWT expiry (15 min access, 7 day refresh) | 📋 | TRD-Section5-Security.md |
| 8.2.5 | MFA for sensitive access | PIN verification for sensitive operations | 📋 | PRD-Section2-F1-Authentication.md |
| 8.3.1 | Strong cryptography for authentication | bcrypt cost 12; RS256 JWT | 📋 | TRD-Section5-Security.md |
| 8.3.4 | Invalid attempts limited | Rate limiting (5/minute for login) | 📋 | PRD-Section2-F1-Authentication.md |
| 8.3.5 | Password complexity | Min 8 chars; upper, lower, number, special | 📋 | PRD Section 5 |
| 8.3.6 | Password history | Last 5 passwords blocked | 📋 | TRD-Section3-DatabaseDesign.md |
| 8.3.7 | First-time password change | force_password_change flag | 📋 | TRD-Section3-DatabaseDesign.md |
| 8.3.9 | Password/PIN different from user ID | Validation rules prevent this | 📋 | PRD Section 5 |
| 8.3.10 | Service accounts have unique credentials | Application-level API keys | 📋 | PRD-Section2-F3-Authorization.md |
| 8.3.10.1 | Service account usage restricted | Scope-limited tokens | 📋 | TRD-Section5-Security.md |
| 8.4.1 | MFA implemented | Password + PIN | 📋 | PRD-Section2-F1-Authentication.md |
| 8.4.2 | MFA for remote access | Required for all API access | 📋 | PRD-Section2-F1-Authentication.md |
| 8.5.1 | Terminated user access revoked | Status = INACTIVE; token blacklist | 📋 | PRD-Section2-F2-UserManagement.md |
| 8.6.1 | Application and system accounts managed | Service accounts tracked; audit logged | 📋 | PRD-Section2-F5-AuditLogging.md |

### 3.3 Requirement 10: Log and Monitor All Access

| Requirement | Control | IAM Implementation | Status | Evidence |
|-------------|---------|-------------------|--------|----------|
| 10.1.1 | Audit log process | Async logging to OpenSearch | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.2.1 | All user access logged | Authentication events logged | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.2.1.1 | Admin actions logged | Admin operations with before/after state | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.2.1.2 | Access to audit logs logged | audit:read permission tracked | 📋 | PRD-Section2-F3-Authorization.md |
| 10.2.1.3 | Failed login attempts logged | LOGIN_FAILED event type | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.2.1.4 | User creation/modification logged | User lifecycle events | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.2.1.5 | Permission changes logged | Role assignment events | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.2.2 | Logs contain required details | user_id, timestamp, action, IP, resource | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.3.1 | Audit trail protected | Immutable logs; access control | 📋 | TRD-Section5-Security.md |
| 10.3.2 | Audit log backup | OpenSearch replication; S3 archive | 📋 | TRD-Section5-Security.md |
| 10.3.3 | Audit log integrity | Checksums; append-only | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.4.1 | Log review process | Audit query API; failed login report | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.4.1.1 | Automated log review | Anomaly detection (future) | 📋 | - |
| 10.5.1 | Retain logs for 12 months | 90 days hot, 7 years archive | 📋 | PRD-Section2-F5-AuditLogging.md |
| 10.6.1 | Time synchronization | UTC timestamps; NTP | 📋 | TRD-Section3-DatabaseDesign.md |

---

## 4. OWASP Top 10 2021 Compliance Checklist

### 4.1 Security Controls

| Rank | Risk | IAM Mitigation | Status | Evidence |
|------|------|----------------|--------|----------|
| A01:2021 | Broken Access Control | RBAC; permission checks; tenant isolation | 📋 | PRD-Section2-F3-Authorization.md |
| A02:2021 | Cryptographic Failures | bcrypt hashing; RS256 JWT; TLS 1.2+; Vault for secrets | 📋 | TRD-Section5-Security.md |
| A03:2021 | Injection | Parameterized queries; input validation | 📋 | TRD-Section4-APIDesign.md |
| A04:2021 | Insecure Design | Security review in SDLC; threat modeling | 📋 | TRD-Section5-Security.md |
| A05:2021 | Security Misconfiguration | Secure defaults; hardened headers; no debug in prod | 📋 | TRD-Section5-Security.md |
| A06:2021 | Vulnerable Components | Dependency scanning (Trivy); regular updates | 📋 | IAM-Implementation-Plan-V2.md |
| A07:2021 | Authentication Failures | MFA (PIN); rate limiting; lockout; secure password policy | 📋 | PRD-Section2-F1-Authentication.md |
| A08:2021 | Software and Data Integrity Failures | JWT signature verification; code signing | 📋 | TRD-Section5-Security.md |
| A09:2021 | Security Logging and Monitoring | Comprehensive audit logging; alerting | 📋 | PRD-Section2-F5-AuditLogging.md |
| A10:2021 | Server-Side Request Forgery (SSRF) | URL validation; allowlist for external calls | 📋 | TRD-Section4-APIDesign.md |

### 4.2 Detailed Security Controls

#### A01: Broken Access Control

| Control | Implementation | Test Method |
|---------|----------------|-------------|
| Deny by default | Permissions required; no implicit access | Attempt access without role |
| Enforce ownership | Users can only access their own data | Cross-user access test |
| Rate limiting | 5 login attempts/minute | Brute force test |
| Disable directory listing | API-only; no static files | Path traversal test |
| JWT validation | Signature, expiry, issuer checked | Token manipulation test |
| Tenant isolation | tenant_id filtering on all queries | Cross-tenant access test |

#### A02: Cryptographic Failures

| Control | Implementation | Verification |
|---------|----------------|--------------|
| Data in transit | TLS 1.2+ required | SSL Labs scan |
| Password storage | bcrypt cost factor 12 | Code review |
| JWT signing | RS256 asymmetric | Token inspection |
| Key rotation | 30-day rotation | Configuration audit |
| Secrets storage | HashiCorp Vault | Infrastructure audit |
| PII encryption | Vault Transit (if configured) | Data audit |

#### A07: Authentication Failures

| Control | Implementation | Test Method |
|---------|----------------|-------------|
| Brute force protection | Rate limiting; account lockout | Repeated login attempts |
| Credential stuffing protection | Rate limit by IP; CAPTCHA (future) | Automated attack test |
| Session management | JWT with short expiry; refresh rotation | Session hijack test |
| Password policy | Complexity requirements; history check | Weak password test |
| MFA | PIN verification | Login without PIN |

---

## 5. GDPR Compliance Checklist

> **Note:** GDPR applies if the IAM system processes personal data of EU residents.

### 5.1 Data Protection Principles

| Principle | Requirement | IAM Implementation | Status |
|-----------|-------------|-------------------|--------|
| Lawfulness | Legal basis for processing | Consent during registration | 📋 |
| Purpose limitation | Data used only for stated purposes | Profile data used only for IAM | 📋 |
| Data minimization | Collect only necessary data | Minimal required fields | 📋 |
| Accuracy | Keep data accurate and up-to-date | Profile update API | 📋 |
| Storage limitation | Don't keep data longer than necessary | Data retention policy | 📋 |
| Integrity and confidentiality | Protect data appropriately | Encryption; access control | 📋 |
| Accountability | Demonstrate compliance | Audit logs; documentation | 📋 |

### 5.2 Data Subject Rights

| Right | Requirement | IAM Implementation | Status |
|-------|-------------|-------------------|--------|
| Right of access | Provide copy of data | GET /users/me/profile | 📋 |
| Right to rectification | Allow correction | PUT /users/me/profile | 📋 |
| Right to erasure | Delete upon request | Soft delete; anonymization | 📋 |
| Right to restrict processing | Limit processing | Status = SUSPENDED | 📋 |
| Right to data portability | Export in standard format | Export API (future) | 📋 |
| Right to object | Stop processing | Deactivation workflow | 📋 |
| Automated decision-making | Human review option | Manual approval workflow | 📋 |

### 5.3 Security Measures (Article 32)

| Measure | Requirement | IAM Implementation | Status |
|---------|-------------|-------------------|--------|
| Pseudonymization | Replace identifiers | UUIDs; no direct PII in logs | 📋 |
| Encryption | Protect data | TLS; bcrypt; Vault Transit | 📋 |
| Confidentiality | Restrict access | RBAC; tenant isolation | 📋 |
| Integrity | Prevent unauthorized modification | Audit logging; version control | 📋 |
| Availability | Ensure data accessible | HA design; backups | 📋 |
| Resilience | Withstand attacks | Rate limiting; security hardening | 📋 |
| Restoration | Recover from incidents | Backup procedures | 📋 |
| Testing | Verify security measures | Penetration testing | 📋 |

---

## 6. Security Testing Checklist

### 6.1 Authentication Testing

| Test Case | Description | Expected Result | Status |
|-----------|-------------|-----------------|--------|
| AUTH-001 | Login with valid credentials | Success, JWT returned | 📋 |
| AUTH-002 | Login with invalid password | Failure, error message | 📋 |
| AUTH-003 | Login with locked account | Failure, account locked error | 📋 |
| AUTH-004 | Login after 5 failed attempts | Account locked | 📋 |
| AUTH-005 | Login without PIN (PIN required) | Failure, PIN required | 📋 |
| AUTH-006 | Token refresh with valid refresh token | New tokens returned | 📋 |
| AUTH-007 | Token refresh with expired refresh token | Failure, re-login required | 📋 |
| AUTH-008 | Access with blacklisted token | Failure, token revoked | 📋 |
| AUTH-009 | Logout invalidates token | Subsequent requests fail | 📋 |
| AUTH-010 | Password reset with valid token | Password changed | 📋 |

### 6.2 Authorization Testing

| Test Case | Description | Expected Result | Status |
|-----------|-------------|-----------------|--------|
| AUTHZ-001 | Access with required permission | Success | 📋 |
| AUTHZ-002 | Access without required permission | 403 Forbidden | 📋 |
| AUTHZ-003 | Cross-tenant data access | 404 Not Found (not 403) | 📋 |
| AUTHZ-004 | Branch-scoped permission check | Correct scoping | 📋 |
| AUTHZ-005 | Permission caching | < 100ms response | 📋 |
| AUTHZ-006 | Role assignment audit | Event logged | 📋 |
| AUTHZ-007 | Admin creates user in own tenant | Success | 📋 |
| AUTHZ-008 | Admin creates user in other tenant | Failure | 📋 |

### 6.3 Input Validation Testing

| Test Case | Description | Expected Result | Status |
|-----------|-------------|-----------------|--------|
| INPUT-001 | SQL injection in email field | Rejected, error message | 📋 |
| INPUT-002 | XSS in profile fields | Sanitized or rejected | 📋 |
| INPUT-003 | Invalid email format | Validation error | 📋 |
| INPUT-004 | Password below minimum length | Validation error | 📋 |
| INPUT-005 | PIN with non-numeric characters | Validation error | 📋 |
| INPUT-006 | UUID injection | Rejected | 📋 |
| INPUT-007 | Path traversal in file upload | Rejected | 📋 |
| INPUT-008 | JSON injection in metadata | Sanitized | 📋 |

### 6.4 Rate Limiting Testing

| Test Case | Description | Expected Result | Status |
|-----------|-------------|-----------------|--------|
| RATE-001 | 5 login attempts in 1 minute | 6th attempt blocked (429) | 📋 |
| RATE-002 | Rate limit headers present | X-RateLimit-* headers | 📋 |
| RATE-003 | Rate limit reset after window | Access restored | 📋 |
| RATE-004 | Different IPs have separate limits | Independent limits | 📋 |
| RATE-005 | Admin endpoints have lower limits | 30/minute | 📋 |

### 6.5 Session Management Testing

| Test Case | Description | Expected Result | Status |
|-----------|-------------|-----------------|--------|
| SESSION-001 | Access token expiry | Token invalid after 15 min | 📋 |
| SESSION-002 | Refresh token expiry | Token invalid after 7 days | 📋 |
| SESSION-003 | Refresh token rotation | Old refresh token invalid | 📋 |
| SESSION-004 | Concurrent sessions allowed | Both sessions work | 📋 |
| SESSION-005 | Logout terminates session | Token blacklisted | 📋 |

---

## 7. Penetration Testing Scope

### 7.1 In-Scope

| Area | Components | Attack Vectors |
|------|------------|----------------|
| Authentication | /auth/* endpoints | Brute force, credential stuffing, session hijacking |
| Authorization | /users/*, /roles/*, /permissions/* | Privilege escalation, IDOR, broken access control |
| Input Validation | All POST/PUT endpoints | SQL injection, XSS, command injection |
| Session Management | JWT handling | Token manipulation, replay attacks |
| Cryptography | Password storage, JWT signing | Weak algorithms, key exposure |
| API Security | All endpoints | Rate limiting bypass, DoS |
| Data Protection | User data, audit logs | Data leakage, unauthorized access |

### 7.2 Out-of-Scope

| Area | Reason |
|------|--------|
| Infrastructure (cloud provider) | Separate assessment |
| Physical security | Not applicable |
| Social engineering | Not in current scope |
| Client applications | Separate applications |

### 7.3 Rules of Engagement

| Rule | Description |
|------|-------------|
| Testing window | [Defined dates in staging] |
| Testing environment | Staging only |
| Rate limit exceptions | Whitelisted IPs |
| Data handling | No real user data; test data only |
| Communication | Slack channel for urgent findings |
| Reporting | Written report within 5 business days |

---

## 8. Compliance Verification Schedule

### 8.1 Pre-Alpha (Week 4)

| Check | Owner | Status |
|-------|-------|--------|
| Code security scan (gosec) | Backend | 📋 |
| Dependency vulnerability scan (Trivy) | DevOps | 📋 |
| OWASP checklist review | Security | 📋 |

### 8.2 Pre-Beta (Week 14)

| Check | Owner | Status |
|-------|-------|--------|
| ISO 27001 A.9 controls review | Security | 📋 |
| PCI DSS Req 8 controls review | Security | 📋 |
| Security scan (comprehensive) | Security | 📋 |
| Penetration testing kickoff | Security | 📋 |

### 8.3 Pre-Production (Week 18)

| Check | Owner | Status |
|-------|-------|--------|
| Full ISO 27001 checklist | Security | 📋 |
| Full PCI DSS checklist | Security | 📋 |
| Penetration test report review | Security | 📋 |
| Remediation verification | Backend | 📋 |
| Final security sign-off | Security Lead | 📋 |

### 8.4 Post-Production (Ongoing)

| Check | Frequency | Owner |
|-------|-----------|-------|
| Dependency vulnerability scan | Weekly | DevOps |
| Security patch review | Weekly | Security |
| Access review | Quarterly | Security |
| Penetration testing | Annually | Security |
| Compliance audit | Annually | Security |

---

## 9. Evidence Collection

### 9.1 Required Documentation

| Document | Purpose | Owner |
|----------|---------|-------|
| Security Architecture | Technical design evidence | Tech Lead |
| Data Flow Diagrams | Data handling evidence | Tech Lead |
| Access Control Policy | Policy evidence | Security |
| Password Policy | Control evidence | Security |
| Audit Log Samples | Logging evidence | Backend |
| Test Reports | Testing evidence | QA |
| Penetration Test Report | Security testing evidence | Security |
| Incident Response Plan | Procedure evidence | Security |

### 9.2 Audit Log Retention

| Log Type | Retention Period | Storage |
|----------|-----------------|---------|
| Authentication events | 7 years | OpenSearch → S3 |
| Authorization events | 7 years | OpenSearch → S3 |
| Admin actions | 7 years | OpenSearch → S3 |
| Security events | 7 years | OpenSearch → S3 |
| System logs | 90 days | CloudWatch |

---

## 10. Remediation Tracking

### 10.1 Finding Categories

| Severity | Response Time | Examples |
|----------|---------------|----------|
| Critical | 24 hours | Active exploitation, data breach |
| High | 72 hours | Authentication bypass, privilege escalation |
| Medium | 2 weeks | Information disclosure, weak encryption |
| Low | 1 month | Best practice deviations |

### 10.2 Remediation Template

```markdown
## Finding: [Title]

**Severity:** [Critical/High/Medium/Low]
**Status:** [Open/In Progress/Resolved/Accepted Risk]
**Found Date:** [Date]
**Target Resolution:** [Date]

### Description
[Detailed description of the finding]

### Impact
[Business and security impact]

### Affected Components
- [Component 1]
- [Component 2]

### Remediation Steps
1. [Step 1]
2. [Step 2]

### Verification
[How to verify the fix]

### Resolution
**Resolved Date:** [Date]
**Resolution Notes:** [Notes]
**Verified By:** [Name]
```

---

## Document Sign-Off

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Security Lead | | | |
| Tech Lead | | | |
| Compliance Officer | | | |
| Project Manager | | | |

---

**End of Compliance Checklist v1.0**

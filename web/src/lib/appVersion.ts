/** Version string shared by the app, the site, and GitHub releases. Tags use a leading v. */
export function displayVersion(version?: string | null): string {
  const value = (version ?? "").trim().replace(/^v/i, "")
  if (!value || value === "unknown" || value === "internal") return ""
  return value
}

/**
 * Subsystem accent usage guide (chips / icons / diagrams only — never full-page fills)
 *
 * Yggdrasil  teal      primary actions, focus, active nav
 * Bifrost    cyan      node mesh / connectivity
 * Norn       violet    orchestration / scheduling
 * Huginn     silver    outbound observation
 * Muninn     indigo    memory / context
 * Mimir      gold      knowledge / recommendations
 * Heimdall   amber     diagnostics / watchtower
 * Gungnir    copper    precision / targeting
 */

export const subsystemAccents = {
  yggdrasil: 'ygg',
  bifrost: 'bifrost',
  norn: 'norn',
  huginn: 'huginn',
  muninn: 'muninn',
  mimir: 'mimir',
  heimdall: 'heimdall',
  gungnir: 'gungnir',
} as const

export type SubsystemId = keyof typeof subsystemAccents

import { describe, expect, it } from 'vitest'
import { displayVersion } from './appVersion'

describe('displayVersion', () => {
  it('matches a GitHub release tag and a build ldflag', () => {
    expect(displayVersion('v0.1.0-alpha.23')).toBe('0.1.0-alpha.23')
    expect(displayVersion('0.1.0-alpha.23')).toBe('0.1.0-alpha.23')
    expect(displayVersion('0.1.0-dev')).toBe('0.1.0-dev')
    expect(displayVersion('unknown')).toBe('')
  })
})

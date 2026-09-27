import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { ChatMarkdown } from './ChatMarkdown'
import { prepareChatMarkdown } from './displayChatText'

const weather =
  'Here are some weather resources for Juneau, AK today: - [EaseWeather](https://www.easeweather.com/juneau/today): Hourly forecast. - [National Weather Service](https://forecast.weather.gov/MapClick.php?lat=58.3&lon=-134.4): Official forecast.'

describe('prepareChatMarkdown', () => {
  it('breaks an inline link list onto separate lines', () => {
    const prepared = prepareChatMarkdown(weather)
    expect(prepared).toContain('\n- [EaseWeather]')
    expect(prepared).toContain('\n- [National Weather Service]')
  })

  it('keeps a fenced code sample as code', () => {
    const sample = 'Example:\n```json\n{"query":"not a tool"}\n```'
    expect(prepareChatMarkdown(sample)).toContain('```json')
  })
})

describe('ChatMarkdown', () => {
  it('renders markdown links instead of the raw syntax', () => {
    render(<ChatMarkdown text={weather} />)
    const link = screen.getByRole('link', { name: 'EaseWeather' })
    expect(link).toHaveAttribute('href', 'https://www.easeweather.com/juneau/today')
    expect(screen.queryByText(/\]\(https/)).toBeNull()
  })

  it('does not turn a javascript URL into a link', () => {
    render(<ChatMarkdown text={'[click](javascript:alert(1))'} />)
    expect(screen.queryByRole('link')).toBeNull()
  })
})
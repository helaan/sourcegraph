import * as React from 'react'
import * as H from 'history'
import { SearchType } from './SearchResults'
import { NavLink } from 'react-router-dom'
import { appendOrReplaceSearchType } from '../helpers'
import { buildSearchURLQuery } from '../../../../shared/src/util/url'
import { constant } from 'lodash'

interface Props {
    location: H.Location
    history: H.History
    type: SearchType
    query: string
}

const typeToProse: Record<Exclude<SearchType, null>, string> = {
    diff: 'Diffs',
    commit: 'Commits',
    symbol: 'Symbols',
    repo: 'Repos',
}

export const SearchResultTabHeader: React.FunctionComponent<Props> = props => {
    const q = appendOrReplaceSearchType(props.query, props.type)
    const builtURLQuery = buildSearchURLQuery(q)

    const isActiveFunc = constant(location.search === `?${builtURLQuery}`)
    const type = props.type
    return (
        <li className="nav-item e2e-search-result-tab">
            <NavLink
                to={{ pathname: '/search', search: builtURLQuery }}
                className={`nav-link e2e-search-result-tab-${props.type}`}
                activeClassName="active e2e-search-result-tab--active"
                isActive={isActiveFunc}
            >
                {type ? typeToProse[type] : 'Code'}
            </NavLink>
        </li>
    )
}

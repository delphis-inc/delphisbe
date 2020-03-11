import React, { useState } from "react";
import { useQuery } from 'react-apollo'
import { discussionListQuery } from './query';
import DiscussionCard, { Props as DiscussionCardProps } from './DiscussionCard'
import { ListDiscussions_listDiscussions } from '../../types/api.d'

export interface Props {}
export type DiscussionType = ListDiscussions_listDiscussions;

function DiscussionList(props: Props) {
    const { loading, error, data } = useQuery(discussionListQuery);

    if (loading) {
        return 'Loading...';
    }
    if (error) {
        return `Error! ${error.message}`;
    }

    console.log(data);
    //const discussions: ListDiscussions_listDiscussions[] = data.listDiscussions;
    const discussionCardComponents: JSX.Element[] = [];

    data.listDiscussions.forEach((d: DiscussionType, idx: number) => {
        discussionCardComponents.push(DiscussionCard({ discussion: d, key: `${idx}` }))
    })

    return (
        <div className="container">
            {discussionCardComponents}
        </div>
    )
}

export default DiscussionList
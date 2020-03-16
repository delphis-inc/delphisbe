import React from 'react';
import { useQuery } from 'react-apollo'
import { useParams } from 'react-router-dom';
import { discussionQuery } from '../DiscussionList/query';
import { GetDiscussionByID_discussion_posts } from '../../types/api.d';
import DiscussionInput from './DiscussionInput';
import DiscussionPost from './DiscussionPost';

export interface Props {}

function Discussion(props: Props): JSX.Element {
    const { id } = useParams();
    const { loading, error, data } = useQuery(discussionQuery, {
        variables: {discussionID: id},
    });

    if (loading) {
        return <div>'Loading...'</div>
    }
    if (error) {
        return <div>`Error! ${error.message}`</div>;
    }

    const posts = data.discussion.posts.map((p: GetDiscussionByID_discussion_posts, idx: number) => {
        return <DiscussionPost post={p} />
    });

    return (
        <div className="discussion-view justify-center">
            {posts}
            <DiscussionInput discussionID={id} />
        </div>
    )
}

export default Discussion
import React from 'react';
import './style.scss';
import { GetDiscussionByID_discussion_posts } from '../../../types/api'

export interface Props {
    post: GetDiscussionByID_discussion_posts
}

function DiscussionPost(props: Props) {
    return (
        <div className="px-4 py-2 discussion-post" key={props.post.id}>
            <div className="author flex mb-3">
                <i className="material-icons md-24">mood</i>
                Yellow Llama
            </div>
            <div className="content">
                {props.post.content}
            </div>
        </div>
    );
}

export default DiscussionPost;
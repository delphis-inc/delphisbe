import React from "react";
import { useQuery } from 'react-apollo'
import { discussionQuery } from '../query'
import { DiscussionType } from '../'
import './style.css'

export interface Props {
    discussion: DiscussionType
    key: string
}

function DiscussionCard(props: Props) {
    console.log(props.discussion);

    return (
        <div className="row d-flex justify-content-center">
            <div className="offset-xl-3 col-md-6 discussion-card" key={props.key}>
                <div className="title">
                    THIS IS A TITLE
                </div>
                <div className="participants">
                    <div className="description">
                        3 verified, 2 anon
                    </div>
                    <img className="discussion-card_profile-image" 
                      src={props.discussion.moderator.userProfile.profileImageURL} alt=""/>
                </div>
            </div>
        </div>
    )
}

export default DiscussionCard
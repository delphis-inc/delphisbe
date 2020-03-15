import React from "react";
import {Link} from 'react-router-dom';
import { DiscussionType } from '../';
import './style.scss';

export interface Props {
    discussion: DiscussionType
    key: string
    isFeatured: boolean
}

function DiscussionCard(props: Props) {
    function getPublicStatus() {
        let iconName = 'lock_open'
        if (false) {
            iconName = 'lock'
        }
        let icon = <i className="material-icons md-24">{iconName}</i>
       
        return (
            <div className="text-xs flex discussion_card-content-status items-center">
                {icon} PUBLIC
            </div>
        );
    }

    return (
        <Link to={`/d/${props.discussion.id}`}>
            <div className="w-full px-4 py-2 bg-gray-500 discussion-card" key={props.key}>
                {getPublicStatus()}
                <div className="title font-extrabold text-2xl pb-2">
                    {props.discussion.title}
                </div>
                <div className="subtitle font_bold text-l pb-4">
                    this is a subtitle describing what the title is all about
                </div>
                <div className="flex justify-between flex-row footer">
                    <>
                        <div className="moderated-by">
                            MODERATED BY
                        </div>
                        <div className="flex moderator">
                            <img className="profile-image" src={props.discussion.moderator.userProfile.profileImageURL} alt=""/>
                            <div className="pl-2 moderator-name">
                                {props.discussion.moderator.userProfile.displayName}
                            </div>
                        </div>
                    </>
                    <div className="flex justify-between scores w-1/4">
                        <div className="fire">
                            <i className="material-icons md-24">fireplace</i>
                            <div className="score">2.2k</div>
                        </div>
                        <div className="participants">
                            <i className="material-icons md-24">people_alt</i>
                            <div className="score">22</div>
                        </div>
                        <div className="viewers">
                            <i className="material-icons md-24">remove_red_eye</i>
                            <div className="score">1.1k</div>
                        </div>
                    </div>
                </div>
            </div>
        </Link>
    )
}

export default DiscussionCard
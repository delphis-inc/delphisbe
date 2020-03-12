import React from "react";
import {Link} from 'react-router-dom';
import { DiscussionType } from '../';
import './style.css';

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
                <div className="discussion-card_title font-extrabold text-2xl pb-2">
                    #title-here
                </div>
                <div className="discussion-card_subtitle font_bold text-l pb-4">
                    this is a subtitle describing what the title is all about
                </div>
                <div className="flex justify-between flex-row discussion-card_footer">
                    <div className="fdiscussion-card_moderator-footer">
                        <div className="discussion-card_moderator-footer_moderated-by">
                            MODERATED BY
                        </div>
                        <div className="flex dicussion-card_moderator-footer_moderator">
                            <img className="discussion-card_moderator-footer_moderator_profile-image" src={props.discussion.moderator.userProfile.profileImageURL} alt=""/>
                            <div className="pl-2 discussion-card_moderator-footer_moderator_moderator-name">
                                {props.discussion.moderator.userProfile.displayName}
                            </div>
                        </div>
                    </div>
                    <div className="flex justify-between discussion-card_footer_scores w-1/3">
                        <div className="discussion-card_footer_scores_fire">
                            <i className="material-icons md-24">fireplace</i>
                            <div className="discussion-card_footer_scores_fire_score">2.2k</div>
                        </div>
                        <div className="discussion-card_footer_scores_participants">
                            <i className="material-icons md-24">people_alt</i>
                            <div className="discussion-card_footer_scores_participants_score">22</div>
                        </div>
                        <div className="discussion-card_footer_scores_viewers">
                            <i className="material-icons md-24">remove_red_eye</i>
                            <div className="discussion-card_footer_scores_viewers_score">1.1k</div>
                        </div>
                    </div>
                </div>
            </div>
        </Link>
    )
}

export default DiscussionCard
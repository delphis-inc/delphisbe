import { gql } from 'apollo-boost';
import DiscussionFragments from '../fragments/discussion';

export default {
    createDiscussion: gql`mutation CreateDiscussion{
        createDiscussion($anonymityType: String, $title: String) {
            ...discussionInfo
        }
    }
    ${DiscussionFragments.discussionInfo}`
}
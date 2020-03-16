import { gql } from 'apollo-boost';
import DiscussionFragments from '../fragments/discussion';

export default {
    createPost: gql`mutation CreatePost($discussionID: ID!, $postContent: String!){
        addPost(discussionID: $discussionID, postContent: $postContent) {
            ...discussionInfo
        }
    }
    ${DiscussionFragments.discussionInfo}`
}